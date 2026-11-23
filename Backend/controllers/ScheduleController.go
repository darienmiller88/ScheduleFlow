package controllers

import (
	"ScheduleFlow/Backend/middlewares"
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/services"
	"ScheduleFlow/Backend/utils"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

type toastData struct {
	ToastClass   string
	ToastMessage string
}

type ScheduleController struct {
	Router            *chi.Mux
	templates         *template.Template
	sessionManager    *scs.SessionManager
	scheduleService   services.ScheduleService
	specialistService services.SpecialistService
	uploadService     services.UploadService
}

func NewScheduleController(
	scheduleService services.ScheduleService,
	specialistService services.SpecialistService,
	uploadService services.UploadService,
	sessionManager *scs.SessionManager,
) ScheduleController {
	sc := ScheduleController{
		scheduleService:   scheduleService,
		specialistService: specialistService,
		uploadService:     uploadService,
		templates:         template.Must(template.ParseGlob("templates/partials/*.html")),
		Router:            chi.NewRouter(),
		sessionManager:    sessionManager,
	}

	sc.registerRoutes()

	return sc
}

func (s *ScheduleController) registerRoutes() {
	// Rate limit to 1 request every 3 seconds per IP address for adding a schedule
	s.Router.With(
		httprate.LimitBy(1, 3 * time.Second, middlewares.ClientIPKey),
		middlewares.RequireAuth(s.sessionManager),
	).Post("/upload-schedule", s.uploadScheduleForSpecialist)

	// Rate limit to 1 request every 3 seconds per IP address for queuing a schedule
	s.Router.With(
		httprate.LimitBy(1, 3 * time.Second, middlewares.ClientIPKey),
		middlewares.RequireAuth(s.sessionManager),
	).Post("/queue-schedule", s.queueScheduleForSpecialist)

	// Rate limit to 1 request every 3 seconds per IP address for queuing a schedule
	s.Router.With(
		httprate.LimitBy(1, 3 * time.Second, middlewares.ClientIPKey),
		middlewares.RequireAuth(s.sessionManager),
	).Delete("/dequeue-schedule", s.dequeueScheduleForSpecialist)

	//defaults to rate limiting of 300 requests per minute per IP address
	s.Router.With(middlewares.RequireAuth(s.sessionManager)).Get("/preview-schedule/{specialistID}", s.previewSchedule)

	//defaults to rate limiting of 300 requests per minute per IP address
	s.Router.With(middlewares.RequireAuth(s.sessionManager)).Get("/get-schedule/{specialistID}", s.getScheduleForSpecialist)
}

func (s *ScheduleController) previewSchedule(res http.ResponseWriter, req *http.Request) {

}

func (s *ScheduleController) getScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {

}

func (s *ScheduleController) dequeueScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {
	userId := s.sessionManager.GetInt(req.Context(), "userID")
	result := s.scheduleService.RemoveSchedule(userId)

	if result.Err != nil {
		utils.SendHtmlError(res, http.StatusBadRequest, result.Err.Error())
		return
	}

	data := toastData{
		ToastClass:   "toast-success",
		ToastMessage: "Schedule will not be sent.",
	}

	if err := s.templates.ExecuteTemplate(res, "toast", data); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, err.Error())
	}
}

func (s *ScheduleController) queueScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {
	userIDString      := strconv.Itoa(s.sessionManager.GetInt(req.Context(), "userID"))
	schedulePath, err := s.uploadService.GetScheduleOnDisk(userIDString)
	
	if err != nil{
		utils.SendHtmlError(res, http.StatusBadRequest, err.Error())
		return
	}
	
	userId   := s.sessionManager.GetInt(req.Context(), "userID")
	schedule := models.Schedule{
		SpecialistID: userId,
		FilePath: schedulePath,
		SendDate: utils.GetSendDateForEmail(),
	}

	result := s.scheduleService.AddSchedule(schedule)

	if result.Err != nil {
		utils.SendHtmlError(res, http.StatusBadRequest, result.Err.Error())
		return
	}

	data := toastData{
		ToastClass:   "toast-success",
		ToastMessage: "Schedule queued for next week!",
	}

	if err := s.templates.ExecuteTemplate(res, "toast", data); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, err.Error())
	}
}

func (s *ScheduleController) uploadScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		utils.SendHtmlError(res, http.StatusBadRequest, err.Error())
		return
	}

	file, _, err := req.FormFile("schedule")

	if err != nil {
		if err == http.ErrMissingFile {
			utils.SendHtmlError(res, http.StatusBadRequest, "No file uploaded")
			return
		}

		utils.SendHtmlError(res, http.StatusBadRequest, err.Error())
		return
	}

	defer file.Close()

	//receive the user id from their session.
	userID       := strconv.Itoa(s.sessionManager.GetInt(req.Context(), "userID"))
	uploadResult := s.uploadService.UploadSchedule(file, userID)

	if uploadResult.Err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, uploadResult.Err.Error())
		return
	}

	data := toastData{
		ToastClass:   "toast-success",
		ToastMessage: "Schedule uploaded successfully!",
	}

	if err := s.templates.ExecuteTemplate(res, "toast", data); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, err.Error())
	}
}

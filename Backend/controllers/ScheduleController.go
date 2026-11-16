package controllers

import (
	"ScheduleFlow/Backend/middlewares"
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
		httprate.LimitBy(1, 3*time.Second, middlewares.ClientIPKey),
		middlewares.RequireAuth(s.sessionManager),
	).Post("/add-schedule", s.AddScheduleForSpecialist)

	// Rate limit to 1 request every 3 seconds per IP address for queuing a schedule
	s.Router.With(
		httprate.LimitBy(1, 3*time.Second, middlewares.ClientIPKey),
		middlewares.RequireAuth(s.sessionManager),
	).Put("/queue-schedule", s.QueueScheduleForSpecialist)

	//defaults to rate limiting of 300 requests per minute per IP address
	s.Router.With(middlewares.RequireAuth(s.sessionManager)).Get("/get-schedule/{specialistID}", s.GetScheduleForSpecialist)
}

func (s *ScheduleController) GetScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {

}

func (s *ScheduleController) QueueScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {

}

func (s *ScheduleController) AddScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {
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

	//receive the user id from their session, and use it to get that specialist.
	userID           := strconv.Itoa(s.sessionManager.GetInt(req.Context(), "userID"))
	uploadResult := s.uploadService.UploadSchedule(file, userID)

	if uploadResult.Err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, uploadResult.Err.Error())
		return
	}

	data := toastData{
		ToastClass:   "toast-success",
		ToastMessage: "Schedule uploaded successfully!",
	}

	// Process the uploaded file...
	if err := s.templates.ExecuteTemplate(res, "toast", data); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, err.Error())
	}
}

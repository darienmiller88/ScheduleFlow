package controllers

import (
	"ScheduleFlow/Backend/middlewares"
	"ScheduleFlow/Backend/services"
	"ScheduleFlow/Backend/utils"
	"fmt"
	"html/template"
	"net/http"
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
	Router    *chi.Mux
	templates *template.Template
	sessionManager *scs.SessionManager
	scheduleService services.ScheduleService
}

func NewScheduleController(scheduleService services.ScheduleService, sessionManager *scs.SessionManager) ScheduleController {
	sc := ScheduleController{
		scheduleService: scheduleService,
		templates:       template.Must(template.ParseGlob("templates/partials/*.html")),
		Router:          chi.NewRouter(),
		sessionManager: sessionManager,
	}

	sc.registerRoutes()

	return sc
}

func (s *ScheduleController) registerRoutes() {
	// Rate limit to 1 request every 2 seconds per IP address for adding a schedule
	s.Router.With(
		httprate.LimitBy(1, 2 * time.Second, middlewares.ClientIPKey),
		middlewares.RequireAuth(s.sessionManager),
	).Post("/add-schedule", s.AddScheduleForSpecialist)

	//defaults to rate limiting of 300 requests per minute per IP address
	s.Router.With(middlewares.RequireAuth(s.sessionManager)).Get("/get-schedule/{specialistID}", s.GetScheduleForSpecialist)
}

func (s *ScheduleController) GetScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {
	
}

func (s *ScheduleController) AddScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseMultipartForm(1 << 20); err != nil {
		utils.SendHtmlError(res, http.StatusBadRequest, err.Error())
		return
	}

	file, header, err := req.FormFile("schedule")

	if err != nil {
		utils.SendHtmlError(res, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(header.Filename)

	defer file.Close()

	data := toastData{
		ToastClass:   "toast-success",
		ToastMessage: "Schedule uploaded successfully!",
	}

	// Process the uploaded file...
	if err := s.templates.ExecuteTemplate(res, "toast", data); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, err.Error())
	}
}
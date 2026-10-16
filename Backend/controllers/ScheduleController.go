package controllers

import (
	"ScheduleFlow/Backend/middlewares"
	"ScheduleFlow/Backend/services"
	"html/template"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

type ScheduleController struct {
	Router    *chi.Mux
	templates *template.Template
	sessionManager *scs.SessionManager
	scheduleService services.ScheduleService
}

func NewScheduleController(scheduleService services.ScheduleService) ScheduleController {
	sc := ScheduleController{
		scheduleService: scheduleService,
		templates:       template.Must(template.ParseGlob("templates/partials/*.html")),
		Router:          chi.NewRouter(),
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
	
}
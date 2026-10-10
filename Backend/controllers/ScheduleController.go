package controllers

import (
	"ScheduleFlow/Backend/services"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ScheduleController struct {
	Router    *chi.Mux
	templates *template.Template

	// scheduleService services.
	scheduleService services.ScheduleService
}

func NewScheduleController(scheduleService services.ScheduleService) ScheduleController {
	sc := ScheduleController{
		scheduleService: scheduleService,
		templates:       template.Must(template.ParseGlob("templates/*.html")),
		Router:          chi.NewRouter(),
	}

	sc.registerRoutes()

	return sc
}

func (s *ScheduleController) registerRoutes() {

}

func (s *ScheduleController) GetScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {
	
}

func (s *ScheduleController) AddScheduleForSpecialist(res http.ResponseWriter, req *http.Request) {
	
}
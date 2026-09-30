package controllers

import (
	"ScheduleFlow/Backend/services"
	"html/template"

	"github.com/go-chi/chi/v5"
)

type ScheduleController struct {
	Router    *chi.Mux
	templates *template.Template
	// scheduleService services.
}

func NewScheduleController() ScheduleController {
	sc := ScheduleController{}

	sc.registerRoutes()

	return sc
}

func (s *ScheduleController) registerRoutes() {

}
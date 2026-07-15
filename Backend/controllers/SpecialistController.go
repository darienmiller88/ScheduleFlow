package controllers

import (
	"ScheduleFlow/Backend/services"
	"html/template"

	"github.com/go-chi/chi/v5"
)

type SpecialistController struct {
	Router            *chi.Mux
	templates         *template.Template
	specialistService *services.SpecialistService
}

func NewSpecialistController(specialistService *services.SpecialistService) *SpecialistController {	
	return &SpecialistController{
		Router:            chi.NewRouter(),
		templates:         template.Must(template.ParseGlob("./templates/partials/*.html")),
		specialistService: specialistService,
	}
}


package controllers

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/services"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SpecialistController struct {
	Router            *chi.Mux
	templates         *template.Template
	specialistService services.SpecialistService
}

func NewSpecialistController(specialistService services.SpecialistService) *SpecialistController {	
	return &SpecialistController{
		Router:            chi.NewRouter(),
		templates:         template.Must(template.ParseGlob("./templates/partials/*.html")),
		specialistService: specialistService,
	}
}

func (s *SpecialistController) registerSpecialistRoutes(){
	s.Router.Post("/", s.addNewSpecialist)
}

func (s *SpecialistController) addNewSpecialist(res http.ResponseWriter, req *http.Request){
	req.Body = http.MaxBytesReader(res, req.Body, 1 << 20) // limit to 1 megabyte

	if err := req.ParseForm(); err != nil{
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	specialist := models.Specialist{
		FirstName: req.FormValue("first-name"),
		LastName:  req.FormValue("last-name"),
		Email:     req.FormValue("email"),
		Password:  req.FormValue("password"),
	}

	s.specialistService.AddNewSpecialist(specialist)
}
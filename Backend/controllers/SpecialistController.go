package controllers

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/services"
	"errors"
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
	sc := &SpecialistController{
		Router:            chi.NewRouter(),
		templates:         template.Must(template.ParseGlob("./templates/partials/*.html")),
		specialistService: specialistService,
	}

	sc.registerSpecialistRoutes()

	return sc
}

func (s *SpecialistController) registerSpecialistRoutes() {
	s.Router.Post("/signup", s.signup)
}

func (s *SpecialistController) signup(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		var maxBytesErr *http.MaxBytesError

		if errors.As(err, &maxBytesErr) {
			http.Error(res, "Request payload too large. Must be 1 mb and under", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	specialist := models.Specialist{
		FirstName: req.FormValue("first-name"),
		LastName:  req.FormValue("last-name"),
		Password:  req.FormValue("password"),
		Email:     req.FormValue("email"),
	}

	result := s.specialistService.AddNewSpecialist(specialist)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	res.WriteHeader(http.StatusOK)
	//add cookie, and redirect to home page.
}

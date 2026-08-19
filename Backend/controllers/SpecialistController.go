package controllers

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/services"
	"ScheduleFlow/Backend/middlewares"

	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type SpecialistController struct {
	Router            *chi.Mux
	templates         *template.Template
	specialistService        services.SpecialistService
	emailVerificationService services.EmailVerificationService
	emailSendService         services.EmailSendService
	sessionManager           *scs.SessionManager
}

func NewSpecialistController(
	specialistService services.SpecialistService, 
	emailVerificationService services.EmailVerificationService,
	emailSendService services.EmailSendService,
	sessionManager *scs.SessionManager,
) *SpecialistController {
	sc := &SpecialistController{
		Router:            chi.NewRouter(),
		templates:         template.Must(template.ParseGlob("./templates/partials/*.html")),
		specialistService: specialistService,
		emailVerificationService: emailVerificationService,
		emailSendService: emailSendService,
		sessionManager: sessionManager,
	}

	sc.registerSpecialistRoutes()

	return sc
}

func (s *SpecialistController) registerSpecialistRoutes() {
	s.Router.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAuth(s.sessionManager))
	})

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

	newEmailVerification, err := models.NewEmailVerification(result.ResultData.ID)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	emailVerificationResult := s.emailVerificationService.AddEmailVerificationEntry(newEmailVerification)

	if emailVerificationResult.Err != nil {
		http.Error(res, emailVerificationResult.Err.Error(), emailVerificationResult.StatusCode)
		return
	}

	if err := s.emailSendService.SendVerificationEmail("darienm931@gmail.com", "Darien", emailVerificationResult.ResultData.Code); err != nil {
		http.Error(res, fmt.Sprintf("Failed to send verification email: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("Email verification was sent to email successfully:")
	res.WriteHeader(http.StatusOK)
	//add cookie, and redirect to home page.
}

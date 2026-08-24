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

	s.Router.With(middlewares.SendBackToHome(s.sessionManager)).Post("/signup", s.signUp)
	s.Router.With(middlewares.SendBackToHome(s.sessionManager)).Post("/signin", s.signIn)
	s.Router.With(middlewares.RequireAuth(s.sessionManager)).Post("/signout", s.signOut)
	s.Router.With(middlewares.RequireVerification(s.sessionManager, s.emailVerificationService)).Post("/verify-email", s.verifyEmailCode)
	s.Router.Post("/resend-verification", s.resendVerification)
}

func (s *SpecialistController) signOut(res http.ResponseWriter, req *http.Request){

}

func (s *SpecialistController) verifyEmailCode(res http.ResponseWriter, req *http.Request){
	if err := req.ParseForm(); err != nil{
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	 req.FormValue("verification_code")
}

func (s *SpecialistController) resendVerification(res http.ResponseWriter, req *http.Request){

}

func (s *SpecialistController) signIn(res http.ResponseWriter, req *http.Request){
	
}

func (s *SpecialistController) signUp(res http.ResponseWriter, req *http.Request) {
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

	// Validate the specialist data
	result := s.specialistService.AddNewSpecialist(specialist)

	if result.Err != nil {
		http.Error(res, result.Err.Error(), result.StatusCode)
		return
	}

	// Create a new email verification entry for the newly registered specialist
	newEmailVerification, err := models.NewEmailVerification(result.ResultData.ID)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the email verification entry to the database
	emailVerificationResult := s.emailVerificationService.AddEmailVerificationEntry(newEmailVerification)

	if emailVerificationResult.Err != nil {
		http.Error(res, emailVerificationResult.Err.Error(), emailVerificationResult.StatusCode)
		return
	}

	// Send the verification email
	if err := s.emailSendService.SendVerificationEmail(specialist.Email, specialist.FirstName, emailVerificationResult.ResultData.Code); err != nil {
		http.Error(res, fmt.Sprintf("Failed to send verification email: %v", err), http.StatusInternalServerError)
		return
	}

	// Renew the session token to prevent session fixation attacks
	if err := s.sessionManager.RenewToken(req.Context()); err != nil {
		http.Error(res, fmt.Sprintf("Failed to renew session token: %v", err), http.StatusInternalServerError)
		return
	}

	// Store the user ID in the session after successful signup
	s.sessionManager.Put(req.Context(), "userID", result.ResultData.ID)

	fmt.Println("Email verification was sent to email successfully")
	http.Redirect(res, req, "/verification", http.StatusSeeOther)
}

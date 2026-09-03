package controllers

import (
	"ScheduleFlow/Backend/middlewares"
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/services"
	"ScheduleFlow/Backend/utils"
	"time"

	"fmt"
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type SpecialistController struct {
	Router         *chi.Mux
	templates      *template.Template
	sessionManager *scs.SessionManager

	//Services to be used by controller
	emailSendService         services.EmailSendService
	specialistService        services.SpecialistService
	emailVerificationService services.EmailVerificationService
}

func NewSpecialistController(
	specialistService services.SpecialistService,
	emailVerificationService services.EmailVerificationService,
	emailSendService services.EmailSendService,
	sessionManager *scs.SessionManager,
) *SpecialistController {
	sc := &SpecialistController{
		Router:                   chi.NewRouter(),
		templates:                template.Must(template.ParseGlob("./templates/partials/*.html")),
		specialistService:        specialistService,
		emailVerificationService: emailVerificationService,
		emailSendService:         emailSendService,
		sessionManager:           sessionManager,
	}

	sc.registerSpecialistRoutes()

	return sc
}

func (s *SpecialistController) registerSpecialistRoutes() {
	// s.Router.Group(func(r chi.Router) {
	// 	r.Use(middlewares.RequireAuth(s.sessionManager))
	// })

	s.Router.With(middlewares.SendBackToHome(s.sessionManager)).Post("/signup", s.signUp)
	s.Router.With(middlewares.SendBackToHome(s.sessionManager)).Post("/signin", s.signIn)
	s.Router.With(middlewares.RequireAuth(s.sessionManager)).Post("/signout", s.signOut)
	s.Router.With(middlewares.RequireAuth(s.sessionManager), middlewares.RequireVerification(s.sessionManager, s.emailVerificationService)).Post("/verify-email", s.verifyEmailCode)
	s.Router.With(middlewares.RequireAuth(s.sessionManager), middlewares.RequireVerification(s.sessionManager, s.emailVerificationService)).Post("/resend-verification", s.resendVerification)
	// s.Router.Post("/resend-verification", s.resendVerification)
}

func (s *SpecialistController) signOut(res http.ResponseWriter, req *http.Request) {
	// Renew the session token to prevent session fixation attacks
	if err := s.sessionManager.RenewToken(req.Context()); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, fmt.Sprintf("Failed to renew session token: %v", err))
		return
	}

	// Remove the user ID from the session to log the user out
	s.sessionManager.Pop(req.Context(), "userId")

	// Redirect to the login page after successful logout
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func (s *SpecialistController) verifyEmailCode(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	verficiationCode := req.FormValue("verification_code")
	userId := s.sessionManager.GetInt(req.Context(), "userId")
	result := s.emailVerificationService.VerifyEmailCode(userId, verficiationCode)

	if result.Err != nil {
		utils.SendHtmlError(res, result.StatusCode, result.Err.Error())
		return
	}

	http.Redirect(res, req, "/home", http.StatusSeeOther)
}

// Will be rate limited to 5 a day
func (s *SpecialistController) resendVerification(res http.ResponseWriter, req *http.Request) {
	userId := s.sessionManager.GetInt(req.Context(), "userID")
	specialistResult := s.specialistService.GetSpecialistById(userId)

	if specialistResult.Err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, specialistResult.Err.Error())
		return
	}

	//Generate new code, and code
	code, codeHash := s.emailVerificationService.GenerateNewEmailCode()

	//Create a new email verification to update the old one 
	emailVerificationResult := s.emailVerificationService.UpdateEmailVerificationEntry(models.EmailVerification{
		CodeHash: string(codeHash),
		ExpiresAt: time.Now().Add(15 * time.Minute),
		SpecialistId: userId,
	})

	if emailVerificationResult.Err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, emailVerificationResult.Err.Error())
		return
	}

	// Send the verification email
	if err := s.emailSendService.SendVerificationEmail(specialistResult.ResultData.Email, specialistResult.ResultData.FirstName, code); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, err.Error())
		return
	}

	_, err := res.Write([]byte(`<p style="color: green; font-weight:bold">Verification code re-sent.</p>`))

	if err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *SpecialistController) signIn(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	email      := req.FormValue("email")
	password   := req.FormValue("password")
	rememberMe := req.FormValue("remember-me") == "on"

	// Set session lifetime to 1 year if "Remember Me" is checked
	if rememberMe {
		s.sessionManager.Lifetime = 365 * (7 * 24 * time.Hour) 
	}

	// Authenticate the specialist using the provided email and password
	result := s.specialistService.AuthenticateSpecialist(email, password)

	if result.Err != nil {
		utils.SendHtmlError(res, result.StatusCode, result.Err.Error())
		return
	}

	// Renew the session token to prevent session fixation attacks
	if err := s.sessionManager.RenewToken(req.Context()); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, fmt.Sprintf("Failed to renew session token: %v", err))
		return
	}

	// Store the user ID in the session after successful login
	s.sessionManager.Put(req.Context(), "userID", s.specialistService.GetSpecialistByEmail(email).ResultData.ID)

	// Redirect to the home page after successful login
	http.Redirect(res, req, "/home", http.StatusSeeOther)
}

func (s *SpecialistController) signUp(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
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

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
	"github.com/go-chi/httprate"
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
	// Rate limit to 1 per 3 seconds. This is to prevent spamming the signout endpoint.
	s.Router.With(
		httprate.LimitBy(1, 3 * time.Second, middlewares.ClientIPKey),
		middlewares.RequireAuth(s.sessionManager), 
		middlewares.CheckVerified(s.sessionManager, s.emailVerificationService),
	).Post("/signout", s.signOut)

	// Rate limit to 1 per 3 seconds. This is to prevent spamming the signin endpoint.
	s.Router.With(
		httprate.LimitBy(1, 3 * time.Second, middlewares.ClientIPKey),
		middlewares.SendBackToHome(s.sessionManager),
	).Post("/signin", s.signIn)

	//Rate limit to 1 per 10 seconds. This is to prevent spamming the signup endpoint and creating multiple accounts.
	s.Router.With(
		httprate.LimitBy(1, 10 * time.Second, middlewares.ClientIPKey),
		middlewares.SendBackToHome(s.sessionManager), 
		middlewares.CheckNotVerified(s.sessionManager, s.emailVerificationService),
	).Post("/signup", s.signUp)

	//Rate limit to 1 per 3 seconds. This is to prevent spamming the email verification code.
	s.Router.With(
		httprate.LimitBy(1, 3 * time.Second, middlewares.ClientIPKey),
		middlewares.RequireAuth(s.sessionManager), 
		middlewares.CheckNotVerified(s.sessionManager, s.emailVerificationService),
	).Post("/verify-email", s.verifyEmailCode)

	//Rate limit to 5 per day. This is to prevent spamming the email verification code.
	s.Router.With(
		httprate.LimitBy(5, 24 * time.Hour, middlewares.ClientIPKey),
		middlewares.RequireAuth(s.sessionManager), 
		middlewares.CheckNotVerified(s.sessionManager, s.emailVerificationService),
		middlewares.PreventResendCode(s.sessionManager, s.emailVerificationService),
	).Put("/resend-verification", s.resendVerification)
}

func (s *SpecialistController) signOut(res http.ResponseWriter, req *http.Request) {
	// Renew the session token to prevent session fixation attacks
	if err := s.sessionManager.RenewToken(req.Context()); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, fmt.Sprintf("Failed to renew session token: %v", err))
		return
	}

	// Remove the user ID from the session to log the user out
	s.sessionManager.Pop(req.Context(), "userID")

	// Redirect to the login page after successful logout
	res.Header().Set("HX-Redirect", "/")
    res.WriteHeader(http.StatusOK)
}

func (s *SpecialistController) verifyEmailCode(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		utils.SendHtmlError(res, http.StatusBadRequest, fmt.Sprintf("Failed to parse form: %v", err))
		return
	}

	verificationCode := req.FormValue("verification_code")
	userId := s.sessionManager.GetInt(req.Context(), "userID")
	result := s.emailVerificationService.VerifyEmailCode(userId, verificationCode)

	if result.Err != nil {
		utils.SendHtmlError(res, result.StatusCode, result.Err.Error())
		return
	}

	// If the verification is successful, delete the email verification entry from the database
	deleteResult := s.emailVerificationService.DeleteEmailVerificationEntry(userId)

	if deleteResult.Err != nil {
		utils.SendHtmlError(res, deleteResult.StatusCode, deleteResult.Err.Error())
		return
	}

	// Redirect to the home page after successful verification
	res.Header().Set("HX-Redirect", "/home")
    res.WriteHeader(http.StatusOK)
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
		utils.SendHtmlError(res, http.StatusBadRequest, fmt.Sprintf("Failed to parse form: %v", err))
		return
	}

	email      := req.FormValue("email")
	password   := req.FormValue("password")
	rememberMe := req.FormValue("remember-me") == "on"
	
	// Authenticate the specialist using the provided email and password
	result := s.specialistService.AuthenticateSpecialist(email, password)
	
	if result.Err != nil {
		utils.SendHtmlError(res, result.StatusCode, result.Err.Error())
		return
	}

	// Check if the specialist has already verified their email
	emailVerificationResult := s.emailVerificationService.GetEmailVerificationEntryByEmail(email)
	
	if emailVerificationResult.StatusCode == http.StatusInternalServerError {
		utils.SendHtmlError(res, http.StatusInternalServerError, emailVerificationResult.Err.Error())
		return
	}

	// Renew the session token to prevent session fixation attacks
	if err := s.sessionManager.RenewToken(req.Context()); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, fmt.Sprintf("Failed to renew session token: %v", err))
		return
	}

	// Set session lifetime to 1 year if "Remember Me" is checked
	if rememberMe {
		s.sessionManager.Lifetime = 365 * (24 * time.Hour) 
	}
	
	//if the specialist has not verified their email, redirect them to the verification page
	if emailVerificationResult.StatusCode == http.StatusOK {
		
		//Give the user a new session with their userID so they can access the verification page
		s.sessionManager.Put(req.Context(), "userID", s.specialistService.GetSpecialistByEmail(email).ResultData.ID)

		res.Header().Set("HX-Redirect", "/verification")
		res.WriteHeader(http.StatusOK)
		return
	}

	// Store the user ID in the session after successful login
	s.sessionManager.Put(req.Context(), "userID", s.specialistService.GetSpecialistByEmail(email).ResultData.ID)

	// Reset to default session lifetime of 1 week after setting the user ID
	if rememberMe{
		s.sessionManager.Lifetime = 7 * (24 * time.Hour) 
	}

	// Set the HX-Redirect header to redirect the user to the home page
	res.Header().Set("HX-Redirect", "/home")
	res.WriteHeader(http.StatusOK)
}

func (s *SpecialistController) signUp(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		utils.SendHtmlError(res, http.StatusBadRequest, fmt.Sprintf("Failed to parse form: %v", err))
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
		utils.SendHtmlError(res, result.StatusCode, result.Err.Error())
		return
	}

	// Create a new email verification entry for the newly registered specialist
	newEmailVerification, err := models.NewEmailVerification(result.ResultData.ID)

	if err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, fmt.Sprintf("Failed to create email verification entry: %v", err))
		return
	}

	// Add the email verification entry to the database
	emailVerificationResult := s.emailVerificationService.AddEmailVerificationEntry(newEmailVerification)

	if emailVerificationResult.Err != nil {
		utils.SendHtmlError(res, emailVerificationResult.StatusCode, emailVerificationResult.Err.Error())
		return
	}

	// Send the verification email
	if err := s.emailSendService.SendVerificationEmail(specialist.Email, specialist.FirstName, emailVerificationResult.ResultData.Code); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, fmt.Sprintf("Failed to send verification email: %v", err))
		return
	}

	// Renew the session token to prevent session fixation attacks
	if err := s.sessionManager.RenewToken(req.Context()); err != nil {
		utils.SendHtmlError(res, http.StatusInternalServerError, fmt.Sprintf("Failed to renew session token: %v", err))
		return
	}

	// Store the user ID in the session after successful signup
	s.sessionManager.Put(req.Context(), "userID", result.ResultData.ID)

	// Set the HX-Redirect header to redirect the user to the verification page
	res.Header().Set("HX-Redirect", "/verification")
	res.WriteHeader(http.StatusOK)
}

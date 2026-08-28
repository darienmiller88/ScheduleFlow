package services

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/repositories"
	"ScheduleFlow/Backend/utils"
	"errors"
	"math/rand/v2"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type EmailVerificationService interface {
	
	//Update an email verification enntry with a new expiry and new code hash
	UpdateEmailVerificationEntry(emailVerification models.EmailVerification) models.Result[models.EmailVerification]
	
	//Add a new email veriication entry
	AddEmailVerificationEntry(emailVerification models.EmailVerification) models.Result[models.EmailVerification]
	
	//Delete a new email verification entry
	DeleteEmailVerificationEntry(specialistId int) models.Result[bool]

	//Get a email verification entry by id
	GetEmailVerificationEntry(specialistId int) models.Result[models.EmailVerification]

	//Verify whether or not the user inputted the correct verification
	VerifyEmailCode(specialistId int, verificationCode string) models.Result[bool]

	//Generate a new email code to verify new account
	GenerateNewEmailCode() string
}

type emailVerificationService struct {
	repo repositories.EmailVerificationRepository
}


func NewEmailVerificationService(repo repositories.EmailVerificationRepository) EmailVerificationService {
	return &emailVerificationService{
		repo: repo,
	}
}

func (e *emailVerificationService) GenerateNewEmailCode() string{
	min, max := 100000, 999999
	code := min + rand.IntN(max-min)

	return strconv.Itoa(code)
}

//
func (e *emailVerificationService) VerifyEmailCode(specialistId int, verificationCode string) models.Result[bool] {
	result := e.repo.GetEmailVerification(specialistId)

	if result.Err != nil{
		return utils.GetResult(result.Err, result.StatusCode, false)
	}

	err := bcrypt.CompareHashAndPassword([]byte(result.ResultData.CodeHash), []byte(verificationCode))

	if err == nil {
		return utils.GetResult(nil, http.StatusOK, true)
	}

	return utils.GetResult(errors.New("Verification code does not match"), http.StatusBadRequest, false)
}

// Method to add a new email verification entry to the database for a specialist
func (e *emailVerificationService) AddEmailVerificationEntry(emailVerification models.EmailVerification) models.Result[models.EmailVerification] {
	return e.repo.AddEmailVerification(emailVerification)
}

// DeleteEmailVerificationEntry implements [EmailVerificationService].
func (e *emailVerificationService) DeleteEmailVerificationEntry(specialistId int) models.Result[bool] {
	return e.repo.DeleteEmailVerification(specialistId)
}

// GetEmailVerificationEntry implements [EmailVerificationService].
func (e *emailVerificationService) GetEmailVerificationEntry(specialistId int) models.Result[models.EmailVerification] {
	return e.repo.GetEmailVerification(specialistId)
}

// UpdateEmailVerificationEntry implements [EmailVerificationService].
func (e *emailVerificationService) UpdateEmailVerificationEntry(emailVerification models.EmailVerification) models.Result[models.EmailVerification] {
	//update the email verification entry in the database using the repository

	return e.repo.UpdateEmailVerification(emailVerification)
}

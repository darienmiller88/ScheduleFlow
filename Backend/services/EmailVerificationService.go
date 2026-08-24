package services

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/repositories"
	"ScheduleFlow/Backend/utils"
	"errors"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type EmailVerificationService interface {
	UpdateEmailVerificationEntry(emailVerification models.EmailVerification) models.Result[models.EmailVerification]
	AddEmailVerificationEntry(emailVerification models.EmailVerification) models.Result[models.EmailVerification]
	DeleteEmailVerificationEntry(specialistId int) models.Result[bool]
	GetEmailVerificationEntry(specialistId int) models.Result[models.EmailVerification]
	VerifyEmailCode(specialistId int, verificationCode int) models.Result[bool]
}

type emailVerificationService struct {
	repo repositories.EmailVerificationRepository
}


func NewEmailVerificationService(repo repositories.EmailVerificationRepository) EmailVerificationService {
	return &emailVerificationService{
		repo: repo,
	}
}

// VerifyEmailCode implements [EmailVerificationService].
func (e *emailVerificationService) VerifyEmailCode(specialistId int, verificationCode int) models.Result[bool] {
	result := e.repo.GetEmailVerification(specialistId)

	if result.Err != nil{
		return utils.GetResult(result.Err, result.StatusCode, false)
	}

	err := bcrypt.CompareHashAndPassword([]byte(result.ResultData.CodeHash), []byte(strconv.Itoa(verificationCode)))

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

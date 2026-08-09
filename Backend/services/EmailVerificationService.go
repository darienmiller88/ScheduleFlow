package services

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/repositories"
)

type EmailVerificationService interface {
	UpdateEmailVerificationEntry(emailVerification models.EmailVerification) models.Result[models.EmailVerification]
	AddEmailVerificationEntry(emailVerification models.EmailVerification) models.Result[models.EmailVerification]
	DeleteEmailVerificationEntry(specialistId int) models.Result[bool]
	GetEmailVerificationEntry(specialistId int) models.Result[models.EmailVerification]
}

type emailVerificationService struct {
	repo repositories.EmailVerificationRepository
}

func NewEmailVerificationService(repo repositories.EmailVerificationRepository) EmailVerificationService {
	return &emailVerificationService{
		repo: repo,
	}
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

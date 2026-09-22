package services

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/repositories"
)

type SentEmailService interface {
	// AddSentEmail adds a new sent email entry to the database for a specialist.
	AddSentEmail(sentEmail models.SentEmail) models.Result[models.SentEmail]

	// GetSentEmailBySpecialistId retrieves all sent email entries for a given specialist ID.
	GetSentEmailBySpecialistId(specialistId int) models.Result[[]models.SentEmail]

	// GetAllSentEmails retrieves all sent email entries in the database.
	GetAllSentEmails() models.Result[[]models.SentEmail]
}

type sentEmailService struct {
	// repository for sent emails
	sentEmailRepository repositories.SentEmailRepository
}

func NewSentEmailService(sentEmailRepository repositories.SentEmailRepository) SentEmailService {
	return &sentEmailService{
		sentEmailRepository: sentEmailRepository,
	}
}

// AddSentEmail implements [SentEmailService].
func (s *sentEmailService) AddSentEmail(sentEmail models.SentEmail) models.Result[models.SentEmail] {
	return s.sentEmailRepository.AddSentEmail(sentEmail)
}

// GetAllSentEmails implements [SentEmailService].
func (s *sentEmailService) GetAllSentEmails() models.Result[[]models.SentEmail] {
	return s.sentEmailRepository.GetAllSentEmails()
}

// GetSentEmailBySpecialistId implements [SentEmailService].
func (s *sentEmailService) GetSentEmailBySpecialistId(specialistId int) models.Result[[]models.SentEmail] {
	return s.sentEmailRepository.GetSentEmailBySpecialistId(specialistId)
}


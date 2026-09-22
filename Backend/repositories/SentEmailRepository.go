package repositories

import (
	"ScheduleFlow/Backend/constants"
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/utils"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type SentEmailRepository interface {
	// AddSentEmail adds a new sent email entry to the database for a specialist.
	AddSentEmail(sentEmail models.SentEmail) models.Result[models.SentEmail]

	// GetSentEmailBySpecialistId retrieves all sent email entries for a given specialist ID.
	GetSentEmailBySpecialistId(specialistId int) models.Result[[]models.SentEmail]

	// GetAllSentEmails retrieves all sent email entries in the database.
	GetAllSentEmails() models.Result[[]models.SentEmail]
}

// Implementation of the SentEmailRepository interface using sql
type sentEmailRepository struct {
	// db connection
	db *sqlx.DB
}


func NewSentEmailRepository(db *sqlx.DB) SentEmailRepository {
	return &sentEmailRepository{
		db: db,
	}
}

// AddSentEmail implements [SentEmailRepository].
func (s *sentEmailRepository) AddSentEmail(sentEmail models.SentEmail) models.Result[models.SentEmail] {
	err := s.db.QueryRow(
		constants.AddSentEmail,
		sentEmail.FilePath,
		sentEmail.SpecialistID,
	).Scan(&sentEmail.ID)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.SentEmail{})
	}

	return utils.GetResult(nil, http.StatusOK, sentEmail)
}

// GetSentEmailBySpecialistId implements [SentEmailRepository].
func (s *sentEmailRepository) GetSentEmailBySpecialistId(specialistId int) models.Result[[]models.SentEmail] {
	sentEmails := []models.SentEmail{}

	err := s.db.Select(&sentEmails, constants.GetSentEmailBySpecialistId, specialistId)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.GetResult(
				fmt.Errorf("specialist %d did not send any emails", specialistId),
				http.StatusNotFound,
				[]models.SentEmail{},
			)
		}

		return utils.GetResult(err, http.StatusInternalServerError, []models.SentEmail{})
	}
	
	return utils.GetResult(nil, http.StatusOK, sentEmails)
}

// GetAllSentEmails implements [SentEmailRepository].
func (s *sentEmailRepository) GetAllSentEmails() models.Result[[]models.SentEmail] {
	sentEmails := []models.SentEmail{}

	err := s.db.Select(&sentEmails, constants.GetAllSentEmails)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, []models.SentEmail{})
	}

	return utils.GetResult(nil, http.StatusOK, sentEmails)
}
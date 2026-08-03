package repositories

import (
	"ScheduleFlow/Backend/constants"
	"ScheduleFlow/Backend/utils"
	"net/http"

	"ScheduleFlow/Backend/models"
	"github.com/jmoiron/sqlx"
)

// Interface for the Email Verifications. Defines the methods that can be used to interact with the
// email verification table in the database.
type EmailVerificationRepository interface {
	UpdateEmailVerification(emailVerification models.EmailVerification) models.Result[models.EmailVerification]
	AddEmailVerification(emailVerification models.EmailVerification) models.Result[models.EmailVerification]
	GetEmailVerification(specialistId int) models.Result[models.EmailVerification]
	DeleteEmailVerification(specialistId int) models.Result[bool]
}

// Implementation of the EmailVerificationRepository interface using sql
type emailVerificationRepository struct {
	db *sqlx.DB
}

func NewEmailVerificationRepository(db *sqlx.DB) EmailVerificationRepository {
	return &emailVerificationRepository{
		db: db,
	}
}

// Method to add a new email verification entry to the database for a specialist
func (e *emailVerificationRepository) AddEmailVerification(emailVerification models.EmailVerification) models.Result[models.EmailVerification] {
	err := e.db.QueryRow(
		constants.AddEmailVerification,
		emailVerification.SpecialistId,
		emailVerification.CodeHash,
		emailVerification.ExpiresAt,
	).Scan(&emailVerification.ID)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.EmailVerification{})
	}

	return utils.GetResult(nil, http.StatusCreated, emailVerification)
}

// Method to delete an email verification entry from the database by specialist ID
func (e *emailVerificationRepository) DeleteEmailVerification(specialistId int) models.Result[bool] {
	_, err := e.db.Exec(constants.DeleteEmailVerification, specialistId)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, false)
	}

	return utils.GetResult(nil, http.StatusOK, true)
}

// Method to retrieve an email verification entry from the database by specialist ID
func (e *emailVerificationRepository) GetEmailVerification(specialistId int) models.Result[models.EmailVerification] {
	var emailVerification models.EmailVerification
	
	err := e.db.Get(&emailVerification, constants.GetEmailVerification, specialistId)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.EmailVerification{})
	}

	return utils.GetResult(nil, http.StatusOK, emailVerification)
}

// Method to update an email verification entry in the database by specialist ID
func (e *emailVerificationRepository) UpdateEmailVerification(emailVerification models.EmailVerification) models.Result[models.EmailVerification] {
	_, err := e.db.Exec(constants.UpdateEmailVerification, emailVerification.CodeHash, emailVerification.ExpiresAt, emailVerification.SpecialistId)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.EmailVerification{})
	}

	return utils.GetResult(nil, http.StatusOK, emailVerification)
}
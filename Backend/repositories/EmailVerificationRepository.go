package repositories

import (
	// "ScheduleFlow/Backend/constants"
	// "ScheduleFlow/Backend/utils"
	// "database/sql"
	// "errors"
	// "fmt"
	// "net/http"
	// "github.com/lib/pq"
	
	"github.com/jmoiron/sqlx"
	"ScheduleFlow/Backend/models"
)

// Interface for the Email Verifications. Defines the methods that can be used to interact with the 
// email verification table in the database.
type EmailVerificationRepository interface {
	UpdateSpecialist(specialist models.Specialist) models.Result[models.Specialist]
	AddEmailVerification(specialist models.Specialist) models.Result[models.Specialist]
	GetSpecialistByEmail(email string) models.Result[models.Specialist]
	GetPasswordByEmail(email string) models.Result[string]
	GetSpecialistById(id int) models.Result[models.Specialist]
	DeleteSpecialist(id int) models.Result[bool]
}

// Implementation of the SpecialistRepository interface using sql
type emailVerificationRepository struct {
	db *sqlx.DB
}

func NewEmailVerificationRepository(db *sqlx.DB) SpecialistRepository {
	return &specialistRepository{
		db: db,
	}
}
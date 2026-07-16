package repositories

import (
	"ScheduleFlow/Backend/constants"
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/utils"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type SpecialistRepository interface {
	UpdateSpecialist(specialist models.Specialist) models.Result[models.Specialist]
	AddSpecialist(specialist models.Specialist)    models.Result[models.Specialist]
	GetSpecialistByEmail(email string)             models.Result[models.Specialist]
	GetPasswordByEmail(email string)               models.Result[string]
	GetSpecialistById(id int)                      models.Result[models.Specialist]
	DeleteSpecialist(id int)                       models.Result[bool]
}

// Implementation of the SpecialistRepository interface using sql
type specialistRepository struct {
	db *sqlx.DB
}

func NewSpecialistRepository(db *sqlx.DB) SpecialistRepository {
	return &specialistRepository{
		db: db,
	}
}

// Adds a new specialists to the DB with the credientials provided in the specialist parameter. Returns the newly created specialist.
func (s *specialistRepository) AddSpecialist(specialist models.Specialist) models.Result[models.Specialist] {
	err := s.db.QueryRow(
		constants.AddSpecialist, 
		specialist.FirstName,
		specialist.LastName, 
		specialist.Email, 
		specialist.Password,
	).Scan(&specialist.ID)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Specialist{})
	}

	return utils.GetResult(nil, http.StatusOK, specialist)
}

// Retrieves the password of a specialist from the DB using the email provided by the specialists. Returns the password as a string.
func (s *specialistRepository) GetPasswordByEmail(email string) models.Result[string] {
	var password string
	err := s.db.QueryRow(constants.GetPasswordByEmail, email).Scan(&password)
	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, "")
	}
	return utils.GetResult(nil, http.StatusOK, password)
}

// Retrieves a specialist from the DB using their email. Returns the specialist as a models.Specialist.
func (s *specialistRepository) GetSpecialistByEmail(email string) models.Result[models.Specialist] {
	var specialist models.Specialist
	err := s.db.QueryRow(constants.GetSpecialistByEmail, email).Scan(
		&specialist.ID,
		&specialist.FirstName,
		&specialist.LastName,
		&specialist.Email,
		&specialist.Password,
	)
	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Specialist{})
	}
	return utils.GetResult(nil, http.StatusOK, specialist)
}

// Retrieves a specialist from the DB using their ID. Returns the specialist as a models.Specialist.
func (s *specialistRepository) GetSpecialistById(id int) models.Result[models.Specialist] {
	var specialist models.Specialist
	err := s.db.QueryRow(constants.GetSpecialistById, id).Scan(
		&specialist.ID,
		&specialist.FirstName,
		&specialist.LastName,
		&specialist.Email,
		&specialist.Password,
	)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Specialist{})
	}
	
	return utils.GetResult(nil, http.StatusOK, specialist)
}	

// Updates a specialist in the DB with the information provided in the specialist parameter. Returns the updated specialist.
func (s *specialistRepository) UpdateSpecialist(specialist models.Specialist) models.Result[models.Specialist] {
	return models.Result[models.Specialist]{}
}	

// Deletes a specialist from the DB using their ID. Returns true if the deletion was successful, false otherwise.
func (s *specialistRepository) DeleteSpecialist(id int) models.Result[bool] {
	return models.Result[bool]{}
}
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
	"github.com/lib/pq"
)

// Interface for the SpecialistRepository. Defines the methods that can be used to interact with the specialists table in the database.
type SpecialistRepository interface {
	UpdateSpecialistDB(specialist models.Specialist) models.Result[models.Specialist]
	AddSpecialistDB(specialist models.Specialist) models.Result[models.Specialist]
	GetSpecialistByEmailDB(email string) models.Result[models.Specialist]
	GetPasswordByEmailDB(email string) models.Result[string]
	GetSpecialistByIdDB(id int) models.Result[models.Specialist]
	DeleteSpecialistDB(id int) models.Result[bool]
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
func (s *specialistRepository) AddSpecialistDB(specialist models.Specialist) models.Result[models.Specialist] {
	err := s.db.QueryRow(
		constants.AddSpecialist,
		specialist.FirstName,
		specialist.LastName,
		specialist.Email,
		specialist.Password,
		specialist.NameAbbrev,
	).Scan(&specialist.ID)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch string(pgErr.Code) {
			case "23505": // unique_violation
				return utils.GetResult(
					fmt.Errorf("specialist email '%s' already exists", specialist.Email),
					http.StatusConflict,
					models.Specialist{},
				)
			}
		}

		//Return 500 for any other general DB error
		return utils.GetResult(err, http.StatusInternalServerError, models.Specialist{})
	}

	return utils.GetResult(nil, http.StatusCreated, specialist)
}

// Retrieves the password of a specialist from the DB using the email provided by the specialists. Returns the password as a string.
func (s *specialistRepository) GetPasswordByEmailDB(email string) models.Result[string] {
	var password string

	err := s.db.Get(&password, constants.GetPasswordByEmail, email)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, "")
	}

	return utils.GetResult(nil, http.StatusOK, password)
}

// Retrieves a specialist from the DB using their email. Returns the specialist as a models.Specialist.
func (s *specialistRepository) GetSpecialistByEmailDB(email string) models.Result[models.Specialist] {
	var specialist models.Specialist

	err := s.db.Get(&specialist, constants.GetSpecialistByEmail, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.GetResult(
				fmt.Errorf("specialist with email %s not found", email),
				http.StatusNotFound,
				models.Specialist{},
			)
		}

		return utils.GetResult(err, http.StatusInternalServerError, models.Specialist{})
	}

	return utils.GetResult(nil, http.StatusOK, specialist)
}

// Retrieves a specialist from the DB using their ID. Returns the specialist as a models.Specialist.
func (s *specialistRepository) GetSpecialistByIdDB(id int) models.Result[models.Specialist] {
	var specialist models.Specialist

	err := s.db.Get(&specialist, constants.GetSpecialistById, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.GetResult(
				fmt.Errorf("specialist with id %d not found", id),
				http.StatusNotFound,
				models.Specialist{},
			)
		}

		return utils.GetResult(err, http.StatusInternalServerError, models.Specialist{})
	}

	return utils.GetResult(nil, http.StatusOK, specialist)
}

// Updates a specialist in the DB with the information provided in the specialist parameter. Returns the updated specialist.
func (s *specialistRepository) UpdateSpecialistDB(specialist models.Specialist) models.Result[models.Specialist] {
	_, err := s.db.Exec(
		constants.UpdateSpecialist,
		specialist.FirstName,
		specialist.LastName,
		specialist.Email,
		specialist.Password,
		specialist.ID,
	)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch string(pgErr.Code) {
			case "23505": // unique_violation
				return utils.GetResult(
					fmt.Errorf("specialist email '%s' already exists", specialist.Email),
					http.StatusConflict,
					models.Specialist{},
				)
			}
		}

		//Return 500 for any other general DB error
		return utils.GetResult(err, http.StatusInternalServerError, models.Specialist{})
	}

	return utils.GetResult(nil, http.StatusOK, specialist)
}

// Deletes a specialist from the DB using their ID. Returns true if the deletion was successful, false otherwise.
func (s *specialistRepository) DeleteSpecialistDB(id int) models.Result[bool] {
	sqlResult, err := s.db.Exec(constants.DeleteSpecialist, id)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, false)
	}

	if rowsAffected, _ := sqlResult.RowsAffected(); rowsAffected == 0 {
		return utils.GetResult(fmt.Errorf("specialist with id %d not found", id), http.StatusNotFound, false)
	}

	return utils.GetResult(nil, http.StatusOK, true)
}

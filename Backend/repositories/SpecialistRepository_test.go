package repositories

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"ScheduleFlow/Backend/constants"
	"ScheduleFlow/Backend/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSpecialistRepo(t *testing.T) (sqlmock.Sqlmock, *specialistRepository) {
	mock_db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(mock_db, "postgres")
	repo := &specialistRepository{
		db: sqlxDB,
	}

	t.Cleanup(func() {
		mock_db.Close()
	})

	return mock, repo
}

func TestAddNewSpecialist_Success(t *testing.T) {
	mock, repo := setupSpecialistRepo(t)
	specialist := models.Specialist{
		FirstName: "fred",
		LastName: "Burger",
		Password: "fyfbhcunde",
		Email: "fred@ucpnyc.org",
	}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	mock.ExpectQuery(regexp.QuoteMeta(constants.AddSpecialist)).
		WithArgs(specialist.FirstName, specialist.LastName, specialist.Email, specialist.Password). 
		WillReturnRows(rows)

	result := repo.AddSpecialist(specialist)

	require.NoError(t, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)

	//Check to make sure the correct id is returned.
	assert.Equal(t, 1, result.ResultData.ID)

	//And ensure the query returned the correct fields in the specialist
	assert.Equal(t, specialist.FirstName, result.ResultData.FirstName)
	assert.Equal(t, specialist.LastName, result.ResultData.LastName)
	assert.Equal(t, specialist.Email, result.ResultData.Email)
	assert.Equal(t, specialist.Password, result.ResultData.Password)
}

// Tests for a specialist beign added with an email that is already taken
func TestAddSpecialist_DuplicateEmail(t *testing.T) {
	mock, repo := setupSpecialistRepo(t)

	specialist := models.Specialist{
		FirstName: "Fred",
		LastName:  "Burger",
		Email:     "fred@ucpnyc.org",
		Password:  "password123",
	}

	pgErr := &pq.Error{
		Code: "23505",
	}

	mock.ExpectQuery(regexp.QuoteMeta(constants.AddSpecialist)).
		WithArgs(specialist.FirstName, specialist.LastName, specialist.Email, specialist.Password).
		WillReturnError(pgErr)

	result := repo.AddSpecialist(specialist)

	require.Error(t, result.Err)

	assert.Equal(t, http.StatusConflict, result.StatusCode)
	assert.Contains(t, result.Err.Error(), "already exists")
	assert.Equal(t, models.Specialist{}, result.ResultData)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddSpecialist_DatabaseFailure(t *testing.T) {
	mock, repo := setupSpecialistRepo(t)

	specialist := models.Specialist{
		FirstName: "Fred",
		LastName:  "Burger",
		Email:     "fred@ucpnyc.org",
		Password:  "password123",
	}

	mock.ExpectQuery(regexp.QuoteMeta(constants.AddSpecialist)).
		WithArgs(specialist.FirstName, specialist.LastName, specialist.Email, specialist.Password).
		WillReturnError(errors.New("database unavailable"))

	result := repo.AddSpecialist(specialist)

	require.Error(t, result.Err)

	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, models.Specialist{}, result.ResultData)

	require.NoError(t, mock.ExpectationsWereMet())
}


/////////////////////////
// UPDATE tests
////////////////////////

func TestUpdateSpecialist_Success(t *testing.T){
	mock, repo   := setupSpecialistRepo(t)
	specialistID := 1
	specialist   := models.Specialist{
		ID: specialistID,
		FirstName: "fred",
		LastName: "Burger",
		Password: "fyfbhcunde",
		Email: "fred@ucpnyc.org",
	}

	mock.ExpectExec(regexp.QuoteMeta(constants.UpdateSpecialist)).
		WithArgs(specialist.FirstName, specialist.LastName, specialist.Email, specialist.Password, specialistID). 
		WillReturnResult(sqlmock.NewResult(int64(specialistID), 1))

	result := repo.UpdateSpecialist(specialist)

	assert.Nil(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, specialist, result.ResultData)
}

func TestUpdateSpecialist_DuplicateEmail(t *testing.T){
	mock, repo   := setupSpecialistRepo(t)
	specialistID := 1
	specialist   := models.Specialist{
		ID: specialistID,
		FirstName: "fred",
		LastName: "Burger",
		Password: "fyfbhcunde",
		Email: "fred@ucpnyc.org",
	}

	pgErr := &pq.Error{
		Code: "23505",
	}

	mock.ExpectExec(regexp.QuoteMeta(constants.UpdateSpecialist)).
		WithArgs(specialist.FirstName, specialist.LastName, specialist.Email, specialist.Password, specialistID). 
		WillReturnError(pgErr)

	result := repo.UpdateSpecialist(specialist)

	assert.NotNil(t, result.Err)
	assert.Equal(t, http.StatusConflict, result.StatusCode)
}

func TestUpdateSpecialist_DatabaseFailure(t *testing.T){
	mock, repo   := setupSpecialistRepo(t)
	specialistID := 1
	specialist   := models.Specialist{
		ID: specialistID,
		FirstName: "fred",
		LastName: "Burger",
		Password: "fyfbhcunde",
		Email: "fred@ucpnyc.org",
	}

	mock.ExpectExec(regexp.QuoteMeta(constants.UpdateSpecialist)).
		WithArgs(specialist.FirstName, specialist.LastName, specialist.Email, specialist.Password, specialistID). 
		WillReturnError(fmt.Errorf("Database error"))

	result := repo.UpdateSpecialist(specialist)

	assert.NotNil(t, result.Err)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}


/////////////////////////
// GET tests
////////////////////////
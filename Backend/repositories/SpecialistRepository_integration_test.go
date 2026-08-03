package repositories

import (
	"ScheduleFlow/Backend/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ///////////////////////
// CREATE/POST tests
// //////////////////////
func TestAddValidPlayer_IntegrationTest_Ok(t *testing.T) {
	specialistRepository := NewSpecialistRepository(db)
	specialist := models.Specialist{
		FirstName: "Darien",
		LastName: "Miller",
		Password: "vybudskm",
		Email: "dr.milr@ucpnyc.org",
	}

	result := specialistRepository.AddSpecialistDB(specialist)

	assert.Equal(t, nil, result.Err)
	assert.Equal(t, http.StatusCreated, result.StatusCode)
	assert.Equal(t, specialist.FirstName, result.ResultData.FirstName)
	assert.Equal(t, specialist.LastName, result.ResultData.LastName)
	assert.Equal(t, specialist.Email, result.ResultData.Email)
	assert.Equal(t, specialist.Password, result.ResultData.Password)
	assert.NotZero(t, result.ResultData.ID)
}

func TestAddDuplicateSpecialist_IntegrationTest_Conflict(t *testing.T) {
	repo := NewSpecialistRepository(db)

	specialist := models.Specialist{
		FirstName: "Darien",
		LastName:  "Miller",
		Password:  "password123",
		Email:     "duplicate@ucpnyc.org",
	}

	// Insert succeeds
	firstResult := repo.AddSpecialistDB(specialist)

	require.NoError(t, firstResult.Err)
	require.Equal(t, http.StatusCreated, firstResult.StatusCode)

	// Insert same email again
	secondResult := repo.AddSpecialistDB(specialist)

	require.Error(t, secondResult.Err)

	assert.Equal(t, http.StatusConflict, secondResult.StatusCode)
	assert.Equal(t, models.Specialist{}, secondResult.ResultData)
	assert.Contains(t, secondResult.Err.Error(), "already exists")
}


/////////////////////////
// READ/GET tests
////////////////////////
func TestGetSpecialistById_IntegrationTest_Success(t *testing.T) {
	repo := NewSpecialistRepository(db)

	//The db is seeded with three specialists, so pick the first one.
	result := repo.GetSpecialistByIdDB(1)

	assert.Nil(t, result.Err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "Darien", result.ResultData.FirstName)
	assert.Equal(t, "Miller", result.ResultData.LastName)
}

func TestGetSpecialistById_IntegrationTest_NotFound(t *testing.T) {
	repo := NewSpecialistRepository(db)

	result := repo.GetSpecialistByIdDB(0)

	assert.NotNil(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
}



/////////////////////////
// UPDATE tests
////////////////////////
func TestUpdateSpecialistById_IntegrationTest_Success(t *testing.T){
	repo := NewSpecialistRepository(db)

	specialist := models.Specialist{
		FirstName: "Darien",
		LastName:  "Miller",
		Password:  "password123",
		Email:     "newemail@ucpnyc.org",
	}

	result := repo.AddSpecialistDB(specialist)

	result.ResultData.FirstName = "Marky"
	result.ResultData.LastName = "Greg"

	updateResult := repo.UpdateSpecialistDB(result.ResultData)

	assert.Nil(t, updateResult.Err)
	assert.Equal(t, http.StatusOK, updateResult.StatusCode)
}

func TestUpdateSpecialistById_IntegrationTest_DuplicateEmail(t *testing.T) {
	repo := NewSpecialistRepository(db)

	// Create first specialist
	first := models.Specialist{
		FirstName: "Darien",
		LastName:  "Miller",
		Password:  "password123",
		Email:     "existing@ucpnyc.org",
	}

	firstResult := repo.AddSpecialistDB(first)
	require.NoError(t, firstResult.Err)
	require.Equal(t, http.StatusCreated, firstResult.StatusCode)

	// Create second specialist
	second := models.Specialist{
		FirstName: "Mark",
		LastName:  "Greg",
		Password:  "password456",
		Email:     "newemail22@ucpnyc.org",
	}

	secondResult := repo.AddSpecialistDB(second)
	require.NoError(t, secondResult.Err)
	require.Equal(t, http.StatusCreated, secondResult.StatusCode)

	// Attempt to update second specialist to use first specialist's email
	secondResult.ResultData.Email = firstResult.ResultData.Email

	updateResult := repo.UpdateSpecialistDB(secondResult.ResultData)

	require.Error(t, updateResult.Err)
	assert.Equal(t, http.StatusConflict, updateResult.StatusCode)
	assert.Equal(t, models.Specialist{}, updateResult.ResultData)
	assert.Contains(t, updateResult.Err.Error(), "already exists")
}



/////////////////////////
// DELETE tests
////////////////////////
func TestDeleteSpecialist_IntegrationTest_Ok(t *testing.T) {
	repo := NewSpecialistRepository(db)

	// Arrange
	specialist := models.Specialist{
		FirstName: "Delete",
		LastName:  "Me",
		Password:  "password123",
		Email:     "deleteme@adaptcn.org",
	}

	createResult := repo.AddSpecialistDB(specialist)

	require.NoError(t, createResult.Err)
	require.Equal(t, http.StatusCreated, createResult.StatusCode)

	// Delete newly added specialist
	deleteResult := repo.DeleteSpecialistDB(createResult.ResultData.ID)

	// Assert
	require.NoError(t, deleteResult.Err)
	assert.Equal(t, http.StatusOK, deleteResult.StatusCode)
	assert.True(t, deleteResult.ResultData)

	//Check to see if the specialist is deleted
	getResult := repo.GetSpecialistByIdDB(createResult.ResultData.ID)

	require.Error(t, getResult.Err)
	assert.Equal(t, http.StatusNotFound, getResult.StatusCode)
}

func TestDeleteSpecialist_IntegrationTest_NotFound(t *testing.T) {
	repo := NewSpecialistRepository(db)

	result := repo.DeleteSpecialistDB(-1)

	require.Error(t, result.Err)
	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.False(t, result.ResultData)
	assert.Contains(t, result.Err.Error(), "specialist with id")
}
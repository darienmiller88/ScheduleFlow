package services

import (
	"fmt"
	"net/http"
	"strings"

	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/repositories"
	"ScheduleFlow/Backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type SpecialistService interface {

	// AddNewSpecialist adds a new specialist to the database. It takes a models.Specialist object as input 
	// and returns a models.Result containing the newly created specialist or an error if the addition fails.
	AddNewSpecialist(specialist models.Specialist) models.Result[models.Specialist]

	// UpdateSpecialist updates the details of an existing specialist in the database. 
	// It takes a models.Specialist object as input and returns a models.Result 
	//  the updated specialist or an error if the update fails.
	UpdateSpecialist(specialist models.Specialist) models.Result[models.Specialist]

	// DeleteSpecialist removes a specialist from the database using their ID.
	DeleteSpecialist(specialistId int) models.Result[bool]

	// GetSpecialistByEmail retrieves a specialist from the database using their email.
	GetSpecialistById(id int) models.Result[models.Specialist]

	// AuthenticateSpecialist checks if the provided email and password match the stored credentials in the database.
	AuthenticateSpecialist(email string, password string) models.Result[bool]

	// GetSpecialistByEmail retrieves a specialist from the database using their email.
	// Returns the specialist as a models.Specialist.
	GetSpecialistByEmail(email string) models.Result[models.Specialist]
}

type specialistService struct {
	specialistRepo repositories.SpecialistRepository
}

// NewSpecialistService creates a new instance of the SpecialistService with the provided SpecialistRepository.
func NewSpecialistService(specialistRepo repositories.SpecialistRepository) SpecialistService {
	return &specialistService{
		specialistRepo: specialistRepo,
	}
}

//Method to authenticate a specialist by checking if the provided email and password match the stored credentials 
// in the database. Returns a boolean indicating whether the authentication was successful or not.
func (s *specialistService) AuthenticateSpecialist(email string, password string) models.Result[bool] {
	specialistResult := s.GetSpecialistByEmail(email)

	// Compare the provided password with the hashed password stored in the database. Do this 
	// first so that email and password are checked at the same time, and the user is not given
	// a hint about which one is incorrect.
	err := bcrypt.CompareHashAndPassword([]byte(specialistResult.ResultData.Password), []byte(password))

	if err != nil {
		return utils.GetResult(fmt.Errorf("email or password is incorrect"), http.StatusUnauthorized, false)
	}

	if specialistResult.Err != nil {

		// If the specialist is not found by email, return a 401 Unauthorized status code
		if specialistResult.StatusCode == http.StatusNotFound {
			return utils.GetResult(fmt.Errorf("email or password is incorrect"), http.StatusUnauthorized, false)
		}

		return utils.GetResult(specialistResult.Err, specialistResult.StatusCode, false)
	}

	return utils.GetResult(nil, http.StatusOK, true)
}

// GetSpecialistByEmail implements [SpecialistService].
func (s *specialistService) GetSpecialistByEmail(email string) models.Result[models.Specialist] {
	return s.specialistRepo.GetSpecialistByEmailDB(email)
}


// GetSpecialistById implements [SpecialistService].
func (s *specialistService) GetSpecialistById(id int) models.Result[models.Specialist] {
	return s.specialistRepo.GetSpecialistByIdDB(id)
}

// UpdateSpecialist implements [SpecialistService].
func (s *specialistService) UpdateSpecialist(specialist models.Specialist) models.Result[models.Specialist] {
	if err := specialist.Validate(); err != nil {
		return utils.GetResult(err, http.StatusUnprocessableEntity, specialist)
	}

	return s.specialistRepo.UpdateSpecialistDB(specialist)
}

func (s *specialistService) AddNewSpecialist(specialist models.Specialist) models.Result[models.Specialist] {

	//Validate the specialist to ensure the first and last names are the appropiate length, the password
	//has the desired length and number of symbols and numbers, and that the email ends in "@ucpnyc.org"
	if err := specialist.Validate(); err != nil {
		return utils.GetResult(err, http.StatusUnprocessableEntity, specialist)
	}

	//hash password, create new row in email verification table with a code (hashed), and expiry (15 minutes)
	// and finally send confirmation email to work email with this code attached.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(specialist.Password), 10)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Specialist{})
	}

	//Set the hashed password to the specialist's password field, and the abbreviated player name
	//to the first letter of the first name and the first letter of the last name.
	specialist.Password = string(hashedPassword)
	specialist.NameAbbrev = s.getPlayerNameInitials(specialist.FirstName, specialist.LastName)

	return s.specialistRepo.AddSpecialistDB(specialist)
}

func (s *specialistService) DeleteSpecialist(specialistId int) models.Result[bool] {
	return s.specialistRepo.DeleteSpecialistDB(specialistId)
}

// Method to combine both initials and return it as: (J)ane (D)oe -> JD
func (s *specialistService) getPlayerNameInitials(firstName string, lastName string) string {
	return strings.ToUpper(string(firstName[0]) + string(lastName[0]))
}

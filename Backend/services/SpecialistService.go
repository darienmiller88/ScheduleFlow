package services

import (
	"net/http"
	"strings"

	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/repositories"
	"ScheduleFlow/Backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type SpecialistService interface {
	AddNewSpecialist(specialist models.Specialist) models.Result[models.Specialist]
	UpdateSpecialist(specialist models.Specialist) models.Result[models.Specialist]
	DeleteSpecialist(specialistId int) models.Result[bool]
	GetSpecialistById(id int) models.Result[models.Specialist]
}

type specialistService struct {
	specialistRepo repositories.SpecialistRepository
}

func NewSpecialistService(specialistRepo repositories.SpecialistRepository) SpecialistService {
	return &specialistService{
		specialistRepo: specialistRepo,
	}
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

package services

import (
	"net/http"
	"strings"

	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/repositories"
	"ScheduleFlow/Backend/utils"

	"golang.org/x/crypto/bcrypt"
	// "github.com/resend/resend-go/v3"
	// "github.com/alexedwards/scs/v2"
)

type SpecialistService interface {
	AddNewSpecialist(specialist models.Specialist) models.Result[models.Specialist]
	DeleteSpecialist(specialistId int)             models.Result[models.Specialist]
}

type specialistService struct {
    repo repositories.SpecialistRepository
}

func NewSpecialistService(repo repositories.SpecialistRepository) SpecialistService {
    return &specialistService{
        repo: repo,
    }
}

func (s *specialistService) AddNewSpecialist(specialist models.Specialist) models.Result[models.Specialist]{
	
	//Validate the specialist to ensure the first and last names are the appropiate length, the password
	//has the desired length and number of symbols and numbers, and that the email ends in "@ucpnyc.org"
	if err := specialist.Validate(); err != nil{
		return utils.GetResult(err, http.StatusUnprocessableEntity, specialist)
	}

	//hash password, create new row in email verification table with a code (hashed), and expiry (15 minutes)
	// and finally send confirmation email to work email with this code attached.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(specialist.Password), 10)

	if err != nil{
		return utils.GetResult(err, http.StatusInternalServerError, models.Specialist{})
	}

	specialist.Password = string(hashedPassword)

    result := s.repo.AddSpecialist(specialist)
    
	return result
}

func (s *specialistService) DeleteSpecialist(specialistId int) models.Result[models.Specialist]{
	return models.Result[models.Specialist]{}
}

//Method to combine both initials and return it as: (J)ane (D)oe -> JD
func (s *specialistService) getPlayerNameInitials(firstName string, lastName string) string {        
		return strings.ToUpper(string(firstName[0]) + string(lastName[0]))
}
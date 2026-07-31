package services

import (
	"net/http"
	"strings"

	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/repositories"
	"ScheduleFlow/Backend/utils"
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


    result := s.repo.AddSpecialist(specialist)

    if result.Err != nil{
		return utils.GetResult(result.Err, result.StatusCode, result.ResultData)
	}
    
	return models.Result[models.Specialist]{}
}

func (s *specialistService) DeleteSpecialist(specialistId int) models.Result[models.Specialist]{
	return models.Result[models.Specialist]{}
}
 
func (s *specialistService) getPlayerNameInitials(playerName string) string {
        
		//Name is validated before insertion, so it SHOULD have exactly 2 parts, ex -> jane doe
		fields := strings.Fields(playerName)

		//Convert the first and last names into runes so the first character could be extracted 
		firstName, lastName := []rune(fields[0]), []rune(fields[1])

		//Extract the first char from the first name and last
		firstNameInitial, lastNameInitial := string(firstName[0]), string(lastName[0])

		//combine both initials and return it as: (J)ane (D)oe -> JD
		return strings.ToUpper(firstNameInitial + lastNameInitial)
}
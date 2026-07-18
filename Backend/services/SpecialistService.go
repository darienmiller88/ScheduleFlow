package services

import (
	"strings"

	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/repositories"
)

type SpecialistService interface {
	AddNewSpecialist(specialist models.Specialist) models.Result[models.Specialist]
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
    
    
    s.repo.AddSpecialist(specialist)
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
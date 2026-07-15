package services

import(
	"ScheduleFlow/Backend/repositories"
)

type SpecialistService interface {
	
}

type specialistService struct {

}

func NewSpecialistService(repo repositories.SpecialistRepository) SpecialistService {
    return &specialistService{
    }
}
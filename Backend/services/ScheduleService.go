package services

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/repositories"
)

type ScheduleService interface {

	// AddSchedule adds a new schedule entry to the database for a specialist.
	AddSchedule(schedule models.Schedule) models.Result[models.Schedule]

	// GetSchedule retrieves the schedule for a specific specialist by their ID.
	GetSchedule(specialistID string) models.Result[models.Schedule]
}

type scheduleService struct {
	scheduleRepository repositories.ScheduleRepository
}

func NewScheduleService(scheduleRepository repositories.ScheduleRepository) ScheduleService {
	return &scheduleService{
		scheduleRepository: scheduleRepository,
	}
}

// GetSchedule implements [ScheduleService].
func (s *scheduleService) GetSchedule(specialistID string) models.Result[models.Schedule] {
	return s.scheduleRepository.GetScheduleForSpecialist(specialistID)
}

// AddSchedule implements [ScheduleService].
func (s *scheduleService) AddSchedule(schedule models.Schedule) models.Result[models.Schedule] {
	return s.scheduleRepository.AddNewSchedule(schedule)
}

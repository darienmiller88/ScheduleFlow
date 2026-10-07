package services

import "ScheduleFlow/Backend/models"

type ScheduleService interface {

	// AddSchedule adds a new schedule entry to the database for a specialist.
	AddSchedule(schedule models.Schedule) models.Result[models.Schedule]

	
}
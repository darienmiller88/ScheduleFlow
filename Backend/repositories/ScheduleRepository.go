package repositories

import (
	"ScheduleFlow/Backend/models"

	"github.com/jmoiron/sqlx"
)

type ScheduleRepository interface {

	//Add a new schedule for a specialist
	AddNewSchedule(schedule models.Schedule) models.Result[models.Schedule]

	//
	Get
}

type scheduleRepository struct {
	db *sqlx.DB
}
package repositories

import (
	"ScheduleFlow/Backend/constants"
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/utils"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type ScheduleRepository interface {

	//Add/update a new schedule for a specialist
	AddNewSchedule(schedule models.Schedule) models.Result[models.Schedule]

	//Get schedule for a specialist
	GetScheduleForSpecialist(specialistID string) models.Result[models.Schedule]
}

type scheduleRepository struct {
	db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) ScheduleRepository {
	return &scheduleRepository{
		db: db,
	}
}

// AddNewSchedule implements [ScheduleRepository].
func (s *scheduleRepository) AddNewSchedule(schedule models.Schedule) models.Result[models.Schedule] {
	err := s.db.QueryRow(
		constants.AddSchedule, 
		schedule.SpecialistID, 
		schedule.FilePath, 
		schedule.SendDate,
	).Scan(&schedule.ID)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, models.Schedule{})
	}

	return utils.GetResult(nil, http.StatusCreated, schedule)
}

// GetScheduleForSpecialist implements [ScheduleRepository].
func (s *scheduleRepository) GetScheduleForSpecialist(specialistID string) models.Result[models.Schedule] {
	schedule := models.Schedule{}
	err      := s.db.Get(&schedule, constants.GetScheduleBySpecialistId, specialistID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.GetResult(fmt.Errorf("schedule not found for specialist %s", specialistID), http.StatusNotFound, models.Schedule{})
		}

		return utils.GetResult(err, http.StatusInternalServerError, models.Schedule{})
	}

	return utils.GetResult(nil, http.StatusOK, schedule)
}


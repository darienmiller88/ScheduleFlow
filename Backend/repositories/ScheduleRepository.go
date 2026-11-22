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
	"github.com/lib/pq"
)

type ScheduleRepository interface {

	//Add a new schedule for a specialist
	AddNewSchedule(schedule models.Schedule) models.Result[models.Schedule]

	//Update the schedule for a specialist
	UpdateSchedule(schedule models.Schedule) models.Result[models.Schedule]

	//Remove, or dequeue a specialists schedule from being sent
	RemoveScheduleForSpecialist(specialistID int) models.Result[bool]

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

// RemoveScheduleForSpecialist implements [ScheduleRepository].
func (s *scheduleRepository) RemoveScheduleForSpecialist(specialistID int) models.Result[bool] {
	_, err := s.db.Exec(constants.DeleteSchedule, specialistID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return utils.GetResult(errors.New("no schedule found for specialist"), http.StatusNotFound, false)
		}

		return utils.GetResult(err, http.StatusInternalServerError, false)
	}

	return utils.GetResult(nil, http.StatusOK, true)
}

// UpdateSchedule implements [ScheduleRepository].
func (s *scheduleRepository) UpdateSchedule(schedule models.Schedule) models.Result[models.Schedule] {
	panic("unimplemented")
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
		if pgErr, ok := err.(*pq.Error); ok {
			switch string(pgErr.Code) {
			case "23505": // unique_violation
				return utils.GetResult(errors.New("schedule is already queued"), http.StatusConflict, models.Schedule{})
			}
		}

		//Return 500 for any other general DB error
		return utils.GetResult(err, http.StatusInternalServerError, models.Schedule{})
	}

	return utils.GetResult(nil, http.StatusCreated, schedule)
}

// GetScheduleForSpecialist implements [ScheduleRepository].
func (s *scheduleRepository) GetScheduleForSpecialist(specialistID string) models.Result[models.Schedule] {
	schedule := models.Schedule{}
	err := s.db.Get(&schedule, constants.GetScheduleBySpecialistId, specialistID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.GetResult(fmt.Errorf("schedule not found for specalist %s", specialistID), http.StatusNotFound, models.Schedule{})
		}

		return utils.GetResult(err, http.StatusInternalServerError, models.Schedule{})
	}

	return utils.GetResult(nil, http.StatusOK, schedule)
}

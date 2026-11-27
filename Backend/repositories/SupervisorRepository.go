package repositories

import "ScheduleFlow/Backend/models"

type SuperVisorRepository interface {

	//Retrieve all supervisors from supervisors table
	GetSuperVisors() models.Result[models.Supervisor]
}
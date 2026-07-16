package repositories

import(
	"ScheduleFlow/Backend/models"
)

type SpecialistRepository interface {
	AddSpecialist(specialist models.Specialist) models.Result[models.Specialist]
	GetPasswordByEmail(email string) models.Result[string]
	GetSpecialistByEmail(email string) models.Result[models.Specialist]
	GetSpecialistById(id int) models.Result[models.Specialist]
	UpdateSpecialist(specialist models.Specialist) models.Result[models.Specialist]
	DeleteSpecialist(id int) models.Result[bool]
}

// Implementation of the SpecialistRepository interface using sql
type specialistRepository struct {

}

func NewSpecialistRepository() SpecialistRepository {
	return &specialistRepository{}
}

func (r *specialistRepository) AddSpecialist(specialist models.Specialist) models.Result[models.Specialist] {
	return models.Result[models.Specialist]{}
}

func (r *specialistRepository) GetPasswordByEmail(email string) models.Result[string] {
	return models.Result[string]{}
}

func (r *specialistRepository) GetSpecialistByEmail(email string) models.Result[models.Specialist] {
	return models.Result[models.Specialist]{}
}

func (r *specialistRepository) GetSpecialistById(id int) models.Result[models.Specialist] {
	return models.Result[models.Specialist]{}
}	

func (r *specialistRepository) UpdateSpecialist(specialist models.Specialist) models.Result[models.Specialist] {
	return models.Result[models.Specialist]{}
}	

func (r *specialistRepository) DeleteSpecialist(id int) models.Result[bool] {
	return models.Result[bool]{}
}
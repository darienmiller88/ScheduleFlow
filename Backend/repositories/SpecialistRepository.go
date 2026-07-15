package repositories

type SpecialistRepository interface {

}

// Implementation of the SpecialistRepository interface using sql
type specialistRepository struct {

}

func NewSpecialistRepository() SpecialistRepository {
	return &specialistRepository{}
}
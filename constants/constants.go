package constants

const (
	AddSpecialist        string = `INSERT INTO specialists (name, email, password) VALUES ($1, $2, $3)`
	GetPasswordByEmail   string = `SELECT password FROM specialists WHERE email = $1`
	GetSpecialistByEmail string = `SELECT * FROM specialists WHERE email = $1`
	GetSpecialistById    string = `SELECT * FROM specialists WHERE id = $1`
	UpdateSpecialist     string = `UPDATE specialists SET name = $1, email = $2, password = $3 WHERE id = $4`
	DeleteSpecialist     string = `DELETE FROM specialists WHERE id = $1`
)
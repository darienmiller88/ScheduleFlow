package constants

const (

	// Query to add a new specialist to the database
	AddSpecialist        string = `INSERT INTO specialists (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id`

	// Query to retrieve a specialist's password from the database by their email
	GetPasswordByEmail   string = `SELECT password FROM specialists WHERE email = $1`

	// Query to retrieve a specialist from the database by their email
	GetSpecialistByEmail string = `SELECT * FROM specialists WHERE email = $1`

	// Query to retrieve a specialist from the database by their ID
	GetSpecialistById    string = `SELECT * FROM specialists WHERE id = $1`

	// Query to update a specialist's information in the database by their ID
	UpdateSpecialist     string = `
		UPDATE specialists 
		SET 
			updated_at = NOW(),
			first_name = $1, 
			last_name = $2,
			email = $3, 
			password = $4 
		WHERE 
			id = $5
	`

	// Query to delete a specialist from the database by their ID if they choose to close their account.
	DeleteSpecialist     string = `DELETE FROM specialists WHERE id = $1`
)
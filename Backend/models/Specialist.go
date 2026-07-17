package models

import "time"

type Specialist struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	//Relevant fields for the Specialist model
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password_hash"`

	//Not a db field, will hold the first letter of the first and last name, ex -> (D)arien (M)iller = DM
	PlayerNameAbbrev string
}
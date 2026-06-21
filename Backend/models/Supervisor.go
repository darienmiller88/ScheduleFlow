package models

import "time"

type Supervisor struct {
	ID        int       `db:"id"         json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	//Relevant fields for the Supervisor model
	FirstName string `db:"first_name"  json:"first_name"`
	LastName  string `db:"last_name"   json:"last_name"`
	Email     string `db:"email"       json:"email"`
}
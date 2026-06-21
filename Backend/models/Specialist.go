package models

import "time"

type Specialist struct {
	ID        int       `db:"id"         json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	//Relevant fields for the Specialist model
	FirstName string    `db:"firstname"  json:"firstname"`
	LastName  string    `db:"lastname"   json:"lastname"`
	Email     string    `db:"email"      json:"email"`
	Password  string    `db:"password"   json:"password"`
}
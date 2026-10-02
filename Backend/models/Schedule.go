package models

import "time"

type Schedule struct {
	ID           int       `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	//Relevant fields for the Schedule model
	SendDate     time.Time `db:"send_date"`
	FilePath     string    `db:"file_path"`
	SpecialistID int       `db:"specialist_id"`
}
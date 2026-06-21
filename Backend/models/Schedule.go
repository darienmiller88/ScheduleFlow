package models

import "time"

type Schedule struct {
	ID           int       `db:"id" json:"id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`

	//Relevant fields for the Schedule model
	SendDate     time.Time `db:"send_date" json:"send_date"`
	FilePath     string    `db:"file_path" json:"file_path"`
	SpecialistID int       `db:"specialist_id" json:"specialist_id"`
}
package models

import "time"

type SentEmail struct {
	ID           int       `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	FilePath     string    `db:"file_path"`
	SpecialistID int       `db:"specialist_id"`	
}
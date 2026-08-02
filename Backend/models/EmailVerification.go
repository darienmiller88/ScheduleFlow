package models

import "time"

type EmailVerification struct{
	ID            int      `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	//Relevant fields
	SpecialistId int           `db:"specialist_id"`
	ExpiresAt    time.Time     `db:"expires_at"`
	CodeHash     string        `db:"code_hash"`
}
package models

import (
	"math/rand/v2"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)
type EmailVerification struct{
	ID            int      `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	//Relevant db fields
	SpecialistId int       `db:"specialist_id"`
	ExpiresAt    time.Time `db:"expires_at"`
	CodeHash     string    `db:"code_hash"`

	//Actual code that will be sent to the specialist's email. This field is not stored in the DB, 
	//but is used to send the code to the specialist.
	Code         string `db:"-"`
}

func NewEmailVerification(specialistId int) (EmailVerification, error) {
	min, max := 1000, 9999
	code := min + rand.IntN(max-min)

	codeString    := strconv.Itoa(code)
	codeHash, err := bcrypt.GenerateFromPassword([]byte(codeString), 10)
	
	if err != nil {
		return EmailVerification{}, err
	}

	return EmailVerification{
		SpecialistId: specialistId,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
		CodeHash:     string(codeHash),
		Code:         codeString,
	}, nil
}
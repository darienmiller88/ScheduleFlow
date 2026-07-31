package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

const (
	// Min length for first and last names.
	minimumLen int = 2

	// Max length for first and last names
	maximumLen int = 20

	// password min length
	minPasswordLen int = 8
)

var (

	//Password needs to have at least one uppercase
	upper = regexp.MustCompile(`[A-Z]`)

	//Password needs to have at least one lowercase
	lower = regexp.MustCompile(`[a-z]`)

	//Password needs to have at least one number
	number = regexp.MustCompile(`[0-9]`)

	//Password needs to have at least one of these symbols
	symbol = regexp.MustCompile(`[!@#$%^&*()_\-+=\[\]{}|\\:;"'<,>.?/]`)

	//Email must be an actual email, and have "@ucpnyc.org" as the domain
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@ucpnyc\.org$`)
)

type Specialist struct {
	ID               int       `db:"id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`

	//Relevant fields for the Specialist model
	FirstName        string `db:"first_name"`
	LastName         string `db:"last_name"`
	Email            string `db:"email"`
	Password         string `db:"password_hash"`

	//will hold the first letter of the first and last name, ex -> (D)arien (M)iller = DM
	PlayerNameAbbrev string `db:"name_abbrev"`	
}

func (s *Specialist) Validate() error {
	if err := s.validateFirstAndLastName(); err != nil {
		return err
	}

	if err := s.validateEmail(); err != nil {
		return err
	}

	if err := s.validatePassword(); err != nil {
		return err
	}

	return nil
}

func (s *Specialist) validatePassword() error {
	if len(s.Password) < minPasswordLen{
		return fmt.Errorf("Password must have at least %d characters close", minPasswordLen)
	}

	if !upper.MatchString(s.Password) {
		return errors.New("Password must have at least one uppercase letter")
	}

	if !lower.MatchString(s.Password) {
		return errors.New("Password must have at least one lowecase letter")
	}

	if !number.MatchString(s.Password) {
		return errors.New("Password must have at least one number")
	}

	if !symbol.MatchString(s.Password) {
		return errors.New(`password must have at least one symbol: !@#$%^&*()_-+={}[]|\\:;"'<,>.?/`)
	}

	return nil
}

func (s *Specialist) validateEmail() error {
	if !emailRegex.MatchString(s.Email) {
		return fmt.Errorf("email must be a valid @ucpnyc.org email address")
	}

	return nil
}

func (s *Specialist) validateFirstAndLastName() error {
	if len(s.FirstName) < minimumLen || len(s.FirstName) > maximumLen{
		return fmt.Errorf("%s must be between %d and %d characters long", s.FirstName, minimumLen, maximumLen)
	}

	if len(s.LastName) < minimumLen || len(s.LastName) > maximumLen {
		return fmt.Errorf("%s is too short, it must be at least %d characters long", s.FirstName, minimumLen)
	}

	if len(s.FirstName) > maximumLen {
		return fmt.Errorf("%s is too long, it must be at most %d characters long", s.FirstName, minimumLen)
	}

	if len(s.LastName) > maximumLen {
		return fmt.Errorf("%s is too long, it must be at most %d characters long", s.FirstName, minimumLen)
	}

	return nil
}

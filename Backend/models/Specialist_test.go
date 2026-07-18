package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpecialist_Validate(t *testing.T) {
	tests := []struct {
		name    string
		model   Specialist
		wantErr bool
	}{
		{
			name: "Valid Specialist",
			model: Specialist{
				FirstName: "Darien",
				LastName:  "Miller",
				Email:     "darien@ucpnyc.org",
				Password:  "Password1!",
			},
			wantErr: false,
		},
		{
			name: "Invalid First Name",
			model: Specialist{
				FirstName: "D",
				LastName:  "Miller",
				Email:     "darien@ucpnyc.org",
				Password:  "Password1!",
			},
			wantErr: true,
		},
		{
			name: "Invalid Last Name",
			model: Specialist{
				FirstName: "Darien",
				LastName:  "M",
				Email:     "darien@ucpnyc.org",
				Password:  "Password1!",
			},
			wantErr: true,
		},
		{
			name: "Invalid Email",
			model: Specialist{
				FirstName: "Darien",
				LastName:  "Miller",
				Email:     "darien@gmail.com",
				Password:  "Password1!",
			},
			wantErr: true,
		},
		{
			name: "Invalid Password",
			model: Specialist{
				FirstName: "Darien",
				LastName:  "Miller",
				Email:     "darien@ucpnyc.org",
				Password:  "password",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSpecialist_validateFirstAndLastName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		wantErr   bool
	}{
		{"Valid", "Darien", "Miller", false},
		{"First Name Too Short", "D", "Miller", true},
		{"Last Name Too Short", "Darien", "M", true},
		{"First Name Too Long", strings.Repeat("A", 21), "Miller", true},
		{"Last Name Too Long", "Darien", strings.Repeat("A", 21), true},
		{"Minimum Boundary", "Da", "Mi", false},
		{"Maximum Boundary", strings.Repeat("A", 20), strings.Repeat("B", 20), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Specialist{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
			}

			err := s.validateFirstAndLastName()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSpecialist_validateEmailDomain(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"Valid", "darien@ucpnyc.org", false},
		{"Missing Domain", "darien", true},
		{"Wrong Domain", "darien@yahoo.com", true},
		{"Empty", "", true},
		{"Only Domain", "@ucpnyc.org", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Specialist{
				Email: tt.email,
			}

			err := s.validateEmail()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSpecialist_validatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid", "Password1!", false},
		{"No Uppercase", "password1!", true},
		{"No Lowercase", "PASSWORD1!", true},
		{"No Number", "Password!", true},
		{"No Symbol", "Password1", true},
		{"Only Uppercase", "PASSWORD", true},
		{"Only Lowercase", "password", true},
		{"Only Numbers", "12345678", true},
		{"Only Symbols", "!@#$%^&*", true},
		{"All Requirements", "Abc123$%", false},
		{"Too Short", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Specialist{
				Password: tt.password,
			}

			err := s.validatePassword()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
package models

import (
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
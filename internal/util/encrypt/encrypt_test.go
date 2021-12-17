package encrypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	type args struct {
		mainPassword  string
		checkPassword string
	}
	tests := []struct {
		name string
		args    args
		success bool
	}{
		{
			name:    "Should validate hash with successful",
			args:    args{mainPassword: "$2a$14$PAzUvqT3jPLpNVKuOvwzCOb4/Vm4DUTEVut6Tbq.Uy7vRfrc1R1aG", checkPassword: "4456"},
			success: true,
		},

		{
			name:    "Should validate hash not successful",
			args:    args{mainPassword: "$2a$14$PAzUvqT3jPLpNVKuOvwzCOb4/Vm4DUTEVut6Tbq.Uy7vRfrc1R1aG", checkPassword: "44s56"},
			success: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testHashPassword := CheckPasswordHash(test.args.mainPassword, test.args.checkPassword)

			assert.Equalf(t, testHashPassword, test.success, test.name)

		})
	}
}

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should generate hash with successful",
			args:    args{password: "123"},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testHashPassword, err := HashPassword(test.args.password)

			if testHashPassword == "" {
				assert.Equalf(t, test.wantErr, true, test.name)
			}

			assert.Nil(t, err, test.name)

			assert.NotNilf(t, testHashPassword, test.name)

		})
	}
}

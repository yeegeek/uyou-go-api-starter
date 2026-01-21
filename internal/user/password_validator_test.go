package user

import (
	"testing"

	"github.com/yeegeek/uyou-go-api-starter/internal/config"
)

func TestPasswordValidator_Validate(t *testing.T) {
	cfg := &config.SecurityConfig{
		PasswordMinLength:        8,
		PasswordRequireUppercase: true,
		PasswordRequireLowercase: true,
		PasswordRequireNumber:    true,
		PasswordRequireSpecial:   true,
	}

	validator := NewPasswordValidator(cfg)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "Password123!",
			wantErr:  false,
		},
		{
			name:     "too short",
			password: "Pass1!",
			wantErr:  true,
		},
		{
			name:     "missing uppercase",
			password: "password123!",
			wantErr:  true,
		},
		{
			name:     "missing lowercase",
			password: "PASSWORD123!",
			wantErr:  true,
		},
		{
			name:     "missing number",
			password: "Password!",
			wantErr:  true,
		},
		{
			name:     "missing special",
			password: "Password123",
			wantErr:  true,
		},
		{
			name:     "all requirements met",
			password: "MyP@ssw0rd",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasswordValidator_GetPasswordRequirements(t *testing.T) {
	cfg := &config.SecurityConfig{
		PasswordMinLength:        8,
		PasswordRequireUppercase: true,
		PasswordRequireLowercase: true,
		PasswordRequireNumber:    true,
		PasswordRequireSpecial:   true,
	}

	validator := NewPasswordValidator(cfg)
	requirements := validator.GetPasswordRequirements()

	if requirements == "" {
		t.Error("GetPasswordRequirements() returned empty string")
	}

	t.Logf("Password requirements: %s", requirements)
}

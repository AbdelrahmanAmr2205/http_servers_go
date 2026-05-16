package auth

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userId := uuid.New()
	validSecret := "correct-secret"
	wrongSecret := "totally-different-secret"
	tests := map[string]struct {
		inputId          uuid.UUID
		inputSecret      string
		expiresIn        time.Duration
		validationSecret string
		want             uuid.UUID
		expectErr        bool
	}{
		"valid_case": {
			inputId:          userId,
			inputSecret:      validSecret,
			expiresIn:        time.Duration(10 * time.Minute),
			validationSecret: validSecret,
			want:             userId,
			expectErr:        false,
		},
		"expired_token_case": {
			inputId:          userId,
			inputSecret:      validSecret,
			expiresIn:        time.Duration(-10 * time.Minute),
			validationSecret: validSecret,
			want:             uuid.Nil,
			expectErr:        true,
		},
		"wrong_secret_validation": {
			inputId:          userId,
			inputSecret:      validSecret,
			expiresIn:        time.Duration(10 * time.Minute),
			validationSecret: wrongSecret,
			want:             uuid.Nil,
			expectErr:        true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			token, err := MakeJWT(tc.inputId, tc.inputSecret, tc.expiresIn)
			if err != nil {
				t.Fatalf("error in MakeJWT(): %v\n", err)
			}
			got, err := ValidateJWT(token, tc.validationSecret)
			if (err != nil) != tc.expectErr {
				t.Errorf("ValidateJWT() error = %v, expectErr = %v", err, tc.expectErr)
			}
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Errorf("ValidateJWT() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password1 := "password1"
	password2 := "password2"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := map[string]struct {
		password      string
		hash          string
		expectErr     bool
		matchPassword bool
	}{
		"correct_password": {
			password:      password1,
			hash:          hash1,
			expectErr:     false,
			matchPassword: true,
		},
		"incorrect_password": {
			password:      password1,
			hash:          hash2,
			expectErr:     false,
			matchPassword: false,
		},
		"empty_password": {
			password:      "",
			hash:          hash1,
			expectErr:     false,
			matchPassword: false,
		},
		"Not_hashed": {
			password:      password1,
			hash:          password1,
			expectErr:     true,
			matchPassword: false,
		},
		"invalid_hash": {
			password:      password1,
			hash:          "invalid_hash",
			expectErr:     true,
			matchPassword: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			match, err := CheckPasswordHash(tc.password, tc.hash)
			if (err != nil) != tc.expectErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tc.expectErr)
			}
			if !tc.expectErr && match != tc.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tc.matchPassword, match)
			}
		})
	}
}

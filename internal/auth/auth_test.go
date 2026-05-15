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

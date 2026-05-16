package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const TokenTypeAccess = "chirpy-access"

func HashPassword(password string) (string, error) {
	hashedPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("couldn't hash password: %w", err)
	}

	return hashedPass, nil
}

func CheckPasswordHash(password, hashedPass string) (bool, error) {
	valid, err := argon2id.ComparePasswordAndHash(password, hashedPass)
	if err != nil {
		return false, err
	}

	return valid, nil
}

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": TokenTypeAccess,
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(expiresIn)),
		"sub": userId.String(),
	})
	s, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return s, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	userIdString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != TokenTypeAccess {
		return uuid.Nil, errors.New("invalid issuer")
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	header, ok := headers["Authorization"]
	if !ok {
		return "", errors.New("JWT Token not in headers")
	}

	token, ok := strings.CutPrefix(header[0], "Bearer ")
	if !ok {
		return "", errors.New("invalid Authorization header format")
	}

	return token, nil
}

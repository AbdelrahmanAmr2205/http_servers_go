package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/AbdelrahmanAmr2205/http_servers_go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	email := strings.TrimSpace(params.Email)
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		respondWithError(w, http.StatusBadRequest, "Invalid email address", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error(), err)
			return
		}

	}

	valid, err := auth.CheckPasswordHash(strings.TrimSpace(params.Password), user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
	}
	if !valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", err)
		return
	}

	var expiresIn time.Duration
	if params.ExpiresInSeconds == 0 {
		expiresIn = time.Hour
	} else {
		expiresIn = time.Duration(params.ExpiresInSeconds * 1_000_000_000)
	}
	if expiresIn > time.Hour {
		expiresIn = time.Hour
	}

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
	}

	w.Header().Add("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, returnVals{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}

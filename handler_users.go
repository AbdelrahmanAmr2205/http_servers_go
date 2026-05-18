package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/AbdelrahmanAmr2205/http_servers_go/internal/auth"
	"github.com/AbdelrahmanAmr2205/http_servers_go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
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

	hashedPass, err := auth.HashPassword(strings.TrimSpace(params.Password))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			respondWithError(w, http.StatusBadRequest, "Email already exists", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerEditUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	if len(params.Email) == 0 || len(params.Password) == 0 {
		respondWithError(w, http.StatusBadRequest, "missing parameters", errors.New("missing parameters"))
		return
	}
	email := strings.TrimSpace(params.Email)
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		respondWithError(w, http.StatusBadRequest, "Invalid email address", err)
		return
	}

	hashedPass, err := auth.HashPassword(strings.TrimSpace(params.Password))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          email,
		HashedPassword: hashedPass,
		ID:             r.Context().Value(UserIdKey).(uuid.UUID),
	})
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			respondWithError(w, http.StatusBadRequest, "Email already exists", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "couldn't update user", err)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	respondWithJSON(w, http.StatusOK, returnVals{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

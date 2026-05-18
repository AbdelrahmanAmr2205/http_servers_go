package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/AbdelrahmanAmr2205/http_servers_go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id format", err)
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			respondWithError(w, http.StatusNotFound, err.Error(), err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "couldn't get chirp", err)
		return
	}

	userId := r.Context().Value(UserIdKey).(uuid.UUID)
	if userId != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Unauthorized", errors.New("Unauthorized"))
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirp.ID,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

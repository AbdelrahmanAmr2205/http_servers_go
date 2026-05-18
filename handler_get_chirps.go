package main

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/AbdelrahmanAmr2205/http_servers_go/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirps := []database.Chirp{}

	authorIdString := r.URL.Query().Get("author_id")

	if authorIdString == "" {
		var err error
		chirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't fetch chirps", err)
			return
		}
	} else {
		authorID, err := uuid.Parse(authorIdString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error(), err)
			return
		}

		chirps, err = cfg.db.GetChirpsByAuthorID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't fetch chirps", err)
			return
		}
	}

	sortingMethod := r.URL.Query().Get("sort")
	if sortingMethod == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	res := []returnVals{}
	for _, chirp := range chirps {
		res = append(res, returnVals{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

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

	respondWithJSON(w, http.StatusOK, returnVals{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

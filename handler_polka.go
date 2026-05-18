package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeChirpyRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	err = cfg.db.UpgradeToChirpyRed(r.Context(), userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			respondWithError(w, http.StatusNotFound, err.Error(), err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

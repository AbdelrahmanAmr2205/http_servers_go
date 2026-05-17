package main

import (
	"net/http"
	"strings"

	"github.com/AbdelrahmanAmr2205/http_servers_go/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error(), err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

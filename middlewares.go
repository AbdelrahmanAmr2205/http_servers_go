package main

import (
	"context"
	"net/http"

	"github.com/AbdelrahmanAmr2205/http_servers_go/internal/auth"
)

const UserIdKey = "userID"

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareAuth(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error(), err)
			return
		}
		userId, err := auth.ValidateJWT(token, cfg.secretKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error(), err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdKey, userId)

		handler(w, r.WithContext(ctx))
	})
}

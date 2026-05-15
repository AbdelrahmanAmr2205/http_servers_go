package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cfg.fileServerHits.Store(0)
	err := cfg.db.ResetDatabase(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset the database", err)
	}
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}

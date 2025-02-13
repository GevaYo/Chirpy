package main

import "net/http"

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Not on dev platform", nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Swap(0)
	cfg.db.DeleteAllUsers(r.Context())
	w.Write([]byte("Hits: 0"))
}

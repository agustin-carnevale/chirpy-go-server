package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) hitsCountHandler(w http.ResponseWriter, req *http.Request) {
	// set Header
	w.Header().Add("Content-Type", "text/html")

	//Status code
	w.WriteHeader(http.StatusOK)

	// Write the response body
	// text := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())

	// HTML template with dynamic visit count
	html := `<!DOCTYPE html>
  <html>
    <body>
      <h1>Welcome, Chirpy Admin</h1>
      <p>Chirpy has been visited %d times!</p>
    </body>
  </html>`

	text := fmt.Sprintf(html, cfg.fileserverHits.Load())
	w.Write([]byte(text))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden,
			"You cannot call this endpoint if not in DEV mode.", nil)
		return
	}

	// reset fileserverHits to 0
	cfg.fileserverHits.Store(0)

	err := cfg.dbQueries.DeleteAllUsers(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Error deleting/reseting users.", err)
		return
	}

	// set Header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	//Status code
	w.WriteHeader(http.StatusOK)

	// Write the response body
	text := "Hits: 0. No users."
	w.Write([]byte(text))
}

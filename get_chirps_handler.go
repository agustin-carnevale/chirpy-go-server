package main

import (
	"net/http"

	"github.com/agustin-carnevale/chirpy-go-server/internal/database"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from DB.", err)
		return
	}

	respondWithJSON(w, http.StatusOK, databaseChirpToChrip(chirps))

}

// Mapping from database.Chirp to main Chirp so I can customize json:keys
func databaseChirpToChrip(dbChirps []database.Chirp) []Chirp {
	var chirps []Chirp
	for _, c := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}
	return chirps
}

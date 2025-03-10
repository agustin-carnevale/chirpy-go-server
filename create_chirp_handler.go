package main

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/agustin-carnevale/chirpy-go-server/internal/database"
	"github.com/google/uuid"
)

const maxChirpLength int = 140

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Chirp body validation
	if len(params.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Body missing in request", nil)
		return
	}
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// Clean body
	params.Body = replaceBadWords(params.Body)

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	// Mapping from database.Chirp to main Chirp so I can customize json:keys
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func replaceBadWords(s string) string {
	badWords := [3]string{"kerfuffle", "sharbert", "fornax"}
	cleanedBody := s

	for _, word := range badWords {
		regex := regexp.MustCompile("(?i)" + word)
		cleanedBody = regex.ReplaceAllString(cleanedBody, "****")
	}
	return cleanedBody
}

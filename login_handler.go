package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/agustin-carnevale/chirpy-go-server/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type returnVal struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Login incorrect", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Login incorrect", err)
		return
	}

	var expiresIn time.Duration
	durationReceived := time.Duration(params.ExpiresInSeconds) * time.Second
	if params.ExpiresInSeconds == 0 || durationReceived > time.Hour {
		expiresIn = time.Hour
	} else {
		expiresIn = durationReceived
	}

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate jwt", err)
		return
	}

	// Mapping from database.User to main User so I can customize json:keys
	respondWithJSON(w, http.StatusOK, returnVal{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jwtToken,
	})

}

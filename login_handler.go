package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/agustin-carnevale/chirpy-go-server/internal/auth"
	"github.com/agustin-carnevale/chirpy-go-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		// ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type returnVal struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
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

	// Acess Token expiration time 1 hour
	jwtTokenExpiresIn := time.Hour

	jwtToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, jwtTokenExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate jwt", err)
		return
	}

	// Refresh Token expiration time 60 days
	refreshTokenExpiresIn := time.Hour * 24 * 60
	refreshTokenExpiresAt := time.Now().Add(refreshTokenExpiresIn)

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token", err)
		return
	}

	// Save refresh token to DB
	_, err = cfg.dbQueries.NewRefreshToken(r.Context(), database.NewRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: refreshTokenExpiresAt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token to DB", err)
		return
	}

	// Mapping from database.User to main User so I can customize json:keys
	respondWithJSON(w, http.StatusOK, returnVal{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed.Bool,
		Token:        jwtToken,
		RefreshToken: refreshTokenString,
	})

}

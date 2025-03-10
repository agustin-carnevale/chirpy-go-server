package main

import (
	"net/http"
	"time"

	"github.com/agustin-carnevale/chirpy-go-server/internal/auth"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	type returnVal struct {
		Token string `json:"token"`
	}

	// Validate refresh token
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization Header", err)
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// userID, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshTokenString)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't get user from DB", err)
	// 	return
	// }

	// Generate new AccessToken with expiration time 1 hour
	jwtTokenExpiresIn := time.Hour

	jwtAccessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, jwtTokenExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate jwt", err)
		return
	}

	// Mapping from database.User to main User so I can customize json:keys
	respondWithJSON(w, http.StatusOK, returnVal{
		Token: jwtAccessToken,
	})

}

func (cfg *apiConfig) revokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	// Validate refresh token
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization Header", err)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) < 5 {
		return "", errors.New("password is too short (at least 5 characters)")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	if len(password) < 5 {
		return errors.New("password is too short")
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	tokenClaims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	// Sign token with secret key
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	// Parse the token with claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(tokenSecret), nil
	})

	// Handle parsing errors
	if err != nil {
		return uuid.Nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	// Check if the token has expired
	if time.Now().After(claims.ExpiresAt.Time) {
		return uuid.Nil, errors.New("token has expired")
	}

	// Parse UUID from the subject field
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID in token")
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authorizationHeader := headers.Get("Authorization")
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header")
	}
	return parts[1], nil
}

func MakeRefreshToken() (string, error) {
	tokenLengthInBytes := 32
	tokenBytes := make([]byte, tokenLengthInBytes)
	_, err := rand.Read(tokenBytes)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return "", err
	}
	tokenString := hex.EncodeToString(tokenBytes)

	return tokenString, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authorizationHeader := headers.Get("Authorization")
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", errors.New("invalid Authorization header")
	}
	return parts[1], nil
}

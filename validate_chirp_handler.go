package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

const maxChirpLength int = 140

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Body missing in request", nil)
		return
	}
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: replaceBadWords(params.Body),
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

// other way of replacing words (split, toLower, join)
// func getCleanedBody(body string, badWords map[string]struct{}) string {
// 	words := strings.Split(body, " ")
// 	for i, word := range words {
// 		loweredWord := strings.ToLower(word)
// 		if _, ok := badWords[loweredWord]; ok {
// 			words[i] = "****"
// 		}
// 	}
// 	cleaned := strings.Join(words, " ")
// 	return cleaned
// }

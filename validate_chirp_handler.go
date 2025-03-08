package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding req params: %s", err)

		type returnVals struct {
			Error string `json:"error"`
		}
		respBody := returnVals{
			Error: "Something went wrong",
		}
		data, err := json.Marshal(respBody)

		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	if len(params.Body) > 140 {
		type returnVals struct {
			Error string `json:"error"`
		}
		respBody := returnVals{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(respBody)

		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}
	respBody := returnVals{
		Valid: true,
	}
	data, err := json.Marshal(respBody)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

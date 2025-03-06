package main

import (
	"fmt"
	"net/http"
)

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	// set Header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	//Status code
	w.WriteHeader(http.StatusOK)
	// Write the response body
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", readinessHandler)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Uppps something went wrong. Server did not start.")
	}

}

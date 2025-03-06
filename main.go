package main

import (
	"fmt"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Uppps something went wrong. Server did not start.")
	}

}

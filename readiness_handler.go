package main

import "net/http"

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	// set Header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	//Status code
	w.WriteHeader(http.StatusOK)
	// Write the response body
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

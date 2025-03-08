package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	// set Header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	//Status code
	w.WriteHeader(http.StatusOK)
	// Write the response body
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) hitsCountHandler(w http.ResponseWriter, req *http.Request) {
	// set Header
	w.Header().Add("Content-Type", "text/html")

	//Status code
	w.WriteHeader(http.StatusOK)

	// Write the response body
	// text := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())

	// HTML template with dynamic visit count
	html := `<!DOCTYPE html>
  <html>
    <body>
      <h1>Welcome, Chirpy Admin</h1>
      <p>Chirpy has been visited %d times!</p>
    </body>
  </html>`

	text := fmt.Sprintf(html, cfg.fileserverHits.Load())
	w.Write([]byte(text))
}

func (cfg *apiConfig) hitsResetHandler(w http.ResponseWriter, req *http.Request) {
	// reset fileserverHits to 0
	cfg.fileserverHits.Store(0)

	// set Header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	//Status code
	w.WriteHeader(http.StatusOK)

	// Write the response body
	text := "Hits: 0"
	w.Write([]byte(text))
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Route Handlers
	mux := http.NewServeMux()

	// File Server
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// API Handlers
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	// Admin  Handlers
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsCountHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.hitsResetHandler)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Uppps something went wrong. Server did not start.")
	}

}

// func middlewareLog(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("%s %s", r.Method, r.URL.Path)
// 		next.ServeHTTP(w, r)
// 	})
// }

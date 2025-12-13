package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/quizzes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listQuizzes(w, r)
		case http.MethodPost:
			createQuiz(w, r)
		case http.MethodOptions:
			// CORS preflight
			enableCors(w)
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// single-quiz handler (by id)
	mux.HandleFunc("/api/quizzes/", func(w http.ResponseWriter, r *http.Request) {
		getQuiz(w, r)
	})
	mux.HandleFunc("/health", health)

	// Wrap with logger
	log.Println("starting backend on :8080")
	if err := http.ListenAndServe(":8080", logger(mux)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

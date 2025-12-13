package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func listQuizzes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		enableCors(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	enableCors(w)
	rows, err := db.Query("SELECT data FROM quizzes ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Printf("query error: %v", err)
		return
	}
	defer rows.Close()
	out := make([]Quiz, 0)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			log.Printf("scan error: %v", err)
			continue
		}
		var q Quiz
		if err := json.Unmarshal(data, &q); err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}
		out = append(out, q)
	}
	if err := rows.Err(); err != nil {
		log.Printf("rows error: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func createQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		enableCors(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	enableCors(w)
	var q Quiz
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	// ensure id and createdAt
	if q.ID == "" {
		q.ID = q.Title + "-" + time.Now().Format("20060102150405.000000")
	}
	if q.CreatedAt.IsZero() {
		q.CreatedAt = time.Now().UTC()
	}
	data, err := json.Marshal(q)
	if err != nil {
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("INSERT INTO quizzes (id, title, data, created_at) VALUES ($1, $2, $3, $4)", q.ID, q.Title, data, q.CreatedAt)
	if err != nil {
		log.Printf("insert error: %v", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(q)
}

func health(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// getQuiz returns a single quiz by ID: GET /api/quizzes/{id}
func getQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		enableCors(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	enableCors(w)
	id := strings.TrimPrefix(r.URL.Path, "/api/quizzes/")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	var data []byte
	err := db.QueryRow("SELECT data FROM quizzes WHERE id=$1", id).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Printf("query error: %v", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	var q Quiz
	if err := json.Unmarshal(data, &q); err != nil {
		log.Printf("unmarshal error: %v", err)
		http.Error(w, "data error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(q)
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

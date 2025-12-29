package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "modernc.org/sqlite"
)

type Entry struct {
	Headword string `json:"headword"`
	HTML     string `json:"html"`
}

func main() {
	db, err := sql.Open("sqlite", "assets/dictionary.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping to verify database opens
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/lookup", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		q = strings.TrimSpace(q)

		if q == "" {
			http.Error(w, "missing q parameter", http.StatusBadRequest)
			return
		}

		var e Entry

		err := db.QueryRow(`
			SELECT headword, html
			FROM entries
			WHERE lower(headword) = lower(?)
			LIMIT 1
		`, q).Scan(&e.Headword, &e.HTML)

		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(e)
	})

	handler := corsMiddleware(mux)

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

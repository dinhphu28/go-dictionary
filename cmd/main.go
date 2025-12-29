package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

type Manifest struct {
	ID        string `json:"id"`
	ShortName string `json:"short_name"`
	FullName  string `json:"full_name"`
	Database  string `json:"database"`
	Version   string `json:"version"`
}

type Dictionary struct {
	Manifest Manifest
	DB       *sql.DB
	Path     string
}

type Entry struct {
	Headword string `json:"headword"`
	HTML     string `json:"html"`
}

type LookupResult struct {
	ID         string  `json:"id"`
	Dictionary string  `json:"dictionary"`
	FullName   string  `json:"full_name"`
	Entries    []Entry `json:"entries"`
}

// ---- global list of loaded dictionaries ----

type GlobalConfig struct {
	Priority []string `json:"priority"`
}

var globalConfig GlobalConfig

// ---- Start of code ----
var dictionaries []Dictionary

// ---- load all dictionaries from resources/ ----

func loadDictionaries(resourceDir string) error {
	err := filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		manifestPath := filepath.Join(path, "manifest.json")
		if _, err := os.Stat(manifestPath); err != nil {
			return nil // not a dictionary folder
		}

		// read manifest
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			return err
		}

		var m Manifest
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}

		if m.Database == "" {
			return errors.New("manifest missing database field: " + manifestPath)
		}

		dbPath := filepath.Join(path, m.Database)

		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			return err
		}

		if err := db.Ping(); err != nil {
			return err
		}

		dictionaries = append(dictionaries, Dictionary{
			Manifest: m,
			DB:       db,
			Path:     path,
		})

		log.Printf("Loaded dictionary: %s (%s)", m.ShortName, dbPath)
		return nil
	})

	return err
}

// ---- lookup exact-case-insensitive in one db ----

func lookupInDB(db *sql.DB, word string) ([]Entry, error) {
	rows, err := db.Query(`
		SELECT headword, html
		FROM entries
		WHERE lower(headword) = lower(?)
	`, word)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Entry

	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.Headword, &e.HTML); err != nil {
			return nil, err
		}
		result = append(result, e)
	}

	return result, nil
}

// ---- HTTP handler ----

func lookupHandler(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		http.Error(w, "missing q parameter", http.StatusBadRequest)
		return
	}

	type resultWrap struct {
		Res LookupResult
		Ok  bool
		Err error
	}

	resultsCh := make(chan resultWrap)
	total := len(dictionaries)

	// launch one goroutine per database
	for _, d := range dictionaries {
		d := d // capture
		go func() {
			entries, err := lookupInDB(d.DB, q)
			if err != nil {
				resultsCh <- resultWrap{Err: err}
				return
			}

			if len(entries) == 0 {
				resultsCh <- resultWrap{Ok: false}
				return
			}

			resultsCh <- resultWrap{
				Ok: true,
				Res: LookupResult{
					ID:         d.Manifest.ID,
					Dictionary: d.Manifest.ShortName,
					FullName:   d.Manifest.FullName,
					Entries:    entries,
				},
			}
		}()
	}

	var results []LookupResult

	// collect responses
	for i := 0; i < total; i++ {
		r := <-resultsCh

		if r.Err != nil {
			log.Println("lookup error:", r.Err)
			continue
		}

		if r.Ok {
			results = append(results, r.Res)
		}
	}

	if len(results) == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// ---- NEW: sort results before response ----
	order := make(map[string]int)

	for i, id := range globalConfig.Priority {
		order[id] = i
	}
	const big = 1_000_000
	sort.Slice(results, func(i, j int) bool {
		oi, okI := order[results[i].ID]
		oj, okJ := order[results[j].ID]

		if !okI {
			oi = big
		}
		if !okJ {
			oj = big
		}

		return oi < oj
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(results)
}

func main() {
	loadConfig("config.json")

	if err := loadDictionaries("resources"); err != nil {
		log.Fatal("failed to load dictionaries:", err)
	}
	applyPriorityOrder()

	log.Printf("Loaded %d dictionaries\n", len(dictionaries))

	http.Handle("/lookup", corsMiddleware(http.HandlerFunc(lookupHandler)))

	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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

func loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &globalConfig)
}

func applyPriorityOrder() {
	order := map[string]int{}
	for i, id := range globalConfig.Priority {
		order[id] = i
	}

	sort.Slice(dictionaries, func(i, j int) bool {
		return order[dictionaries[i].Manifest.ID] < order[dictionaries[j].Manifest.ID]
	})
}

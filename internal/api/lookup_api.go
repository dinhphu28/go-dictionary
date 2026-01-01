// Package api HTTP handler
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"dinhphu28.com/dictionary/internal/config"
	"dinhphu28.com/dictionary/internal/database"
	"dinhphu28.com/dictionary/internal/lookup"
)

func LookupHandler(
	dictionaries []database.Dictionary,
	globalConfig config.GlobalConfig,
) http.HandlerFunc {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		// NOTE: for performance measurement
		start := time.Now()
		defer func() {
			log.Println("lookup total:", time.Since(start))
		}()

		q := strings.TrimSpace(r.URL.Query().Get("q"))
		if q == "" {
			http.Error(w, "missing q parameter", http.StatusBadRequest)
			return
		}

		results := lookup.LookupAllDictionaries(dictionaries, q)

		if len(results) == 0 {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		sortResultsByPriority(results, globalConfig.Priority)

		writeJSONResponse(w, results)
	}
}

func sortResultsByPriority(
	results []lookup.LookupResult,
	priority []string,
) {
	order := make(map[string]int)

	for i, id := range priority {
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
}

func writeJSONResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

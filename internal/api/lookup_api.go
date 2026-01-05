// Package api HTTP handler
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"dinhphu28.com/dictionary/internal/config"
	"dinhphu28.com/dictionary/internal/database"
	"dinhphu28.com/dictionary/internal/lookup"
)

type LookupHandler struct {
	dictionaries []database.Dictionary
	globalConfig config.GlobalConfig
}

func NewLookupHandler(
	dictionaries []database.Dictionary,
	globalConfig config.GlobalConfig,
) *LookupHandler {
	return &LookupHandler{
		dictionaries: dictionaries,
		globalConfig: globalConfig,
	}
}

func (
	lookupHandler *LookupHandler,
) ServeHTTP(
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

	results := lookup.LookupAllDictionariesAndSort(
		lookupHandler.dictionaries,
		q,
		lookupHandler.globalConfig)

	if len(results) == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	writeJSONResponse(w, results)
}

func writeJSONResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

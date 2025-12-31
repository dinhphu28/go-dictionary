// Package api HTTP handler
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"dinhphu28.com/dictionary/internal"
)

func LookupHandler(
	dictionaries []internal.Dictionary,
	globalConfig internal.GlobalConfig,
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

		type resultWrap struct {
			Res internal.LookupResult
			Ok  bool
			Err error
		}

		resultsCh := make(chan resultWrap)
		total := len(dictionaries)

		// launch one goroutine per database
		for _, d := range dictionaries {
			d := d // capture
			go func() {
				entries, err := internal.LookupInDB(d.DB, q)
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
					Res: internal.LookupResult{
						ID:         d.Manifest.ID,
						Dictionary: d.Manifest.ShortName,
						FullName:   d.Manifest.FullName,
						Entries:    entries,
					},
				}
			}()
		}

		var results []internal.LookupResult

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
}

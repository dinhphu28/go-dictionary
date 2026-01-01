package lookup

import (
	"log"
	"sort"

	"dinhphu28.com/dictionary/internal/config"
	"dinhphu28.com/dictionary/internal/database"
)

type resultWrap struct {
	Res LookupResult
	Ok  bool
	Err error
}

func LookupAllDictionariesAndSort(
	dictionaries []database.Dictionary,
	q string,
	globalConfig config.GlobalConfig,
) []LookupResult {
	lookupResults := lookupAllDictionaries(dictionaries, q)
	sortResultsByPriority(lookupResults, globalConfig.Priority)
	return lookupResults
}

func lookupAllDictionaries(
	dictionaries []database.Dictionary,
	q string,
) []LookupResult {
	resultsCh := make(chan resultWrap)
	total := len(dictionaries)

	// launch one goroutine per database
	for _, d := range dictionaries {
		d := d // capture
		go runLookup(d, q, resultsCh)
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
	return results
}

func sortResultsByPriority(
	results []LookupResult,
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

func runLookup(
	d database.Dictionary,
	q string,
	resultsCh chan<- resultWrap,
) {
	entries, err := database.LookupInDB(d.DB, q)
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
}

package lookup

import (
	"log"

	"dinhphu28.com/dictionary/internal/database"
)

type resultWrap struct {
	Res LookupResult
	Ok  bool
	Err error
}

func LookupAllDictionaries(
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

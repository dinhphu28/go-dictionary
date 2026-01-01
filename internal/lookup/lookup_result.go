// Package internal manifest structure to hold application metadata
package lookup

import "dinhphu28.com/dictionary/internal/database"

type LookupResult struct {
	ID         string           `json:"id"`
	Dictionary string           `json:"dictionary"`
	FullName   string           `json:"full_name"`
	Entries    []database.Entry `json:"entries"`
}

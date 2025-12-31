// Package internal manifest structure to hold application metadata
package internal

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

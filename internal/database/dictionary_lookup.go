// Package internal lookup exact-case-insensitive in one db
package database

import (
	"database/sql"
	"log"
	"time"
)

func LookupInDB(db *sql.DB, word string) ([]Entry, error) {
	// NOTE: for performance measurement
	t := time.Now()
	defer func() {
		log.Println("single db lookup:", time.Since(t))
	}()

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

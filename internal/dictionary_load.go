// Package internal load all dictionaries from resources/
package internal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"dinhphu28.com/dictionary/internal/database"
)

var dictionaries []database.Dictionary

func LoadDictionaries(resourceDir string) error {
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

		var m database.Manifest
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

		dictionaries = append(dictionaries, database.Dictionary{
			Manifest: m,
			DB:       db,
			Path:     path,
		})

		log.Printf("Loaded dictionary: %s (%s)", m.ShortName, dbPath)
		return nil
	})

	return err
}

func GetDictionaries() []database.Dictionary {
	return dictionaries
}

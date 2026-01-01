// Package internal manifest structure to hold application metadata
package internal

import (
	"database/sql"

	"dinhphu28.com/dictionary/internal/database"
)

type Dictionary struct {
	Manifest database.Manifest
	DB       *sql.DB
	Path     string
}

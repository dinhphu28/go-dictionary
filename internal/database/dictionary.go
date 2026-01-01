// Package internal manifest structure to hold application metadata
package database

import (
	"database/sql"
)

type Dictionary struct {
	Manifest Manifest
	DB       *sql.DB
	Path     string
}

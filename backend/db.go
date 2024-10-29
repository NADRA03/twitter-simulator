package twitter

import (
	"database/sql"
	"log"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitializeDB initializes the database connection
func InitializeDB() {
    var err error
    db, err = sql.Open("sqlite3", "backend/forum.db") // Make sure this path is correct
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }

    // Log the database path for verification
    dbPath, err := filepath.Abs("./forum.db")
    if err != nil {
        log.Fatalf("Failed to get absolute path of database: %v", err)
    }
    log.Printf("Using database at: %s", dbPath)

    if err = db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	log.Println("Returned db.") 
	return db
}
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

    schema := getSchema()
	// Create all tables
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("error creating tables: %v", err)
	}

    if err = db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }
}

func getSchema() string {
	return `

		CREATE TABLE IF NOT EXISTS Posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS Categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	-- PostCategories table
	CREATE TABLE IF NOT EXISTS PostCategories (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES Posts(id),
		FOREIGN KEY (category_id) REFERENCES Categories(id)
	);

	CREATE TABLE IF NOT EXISTS Comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		username TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES Posts(id),
		FOREIGN KEY (user_id) REFERENCES Users(id)
	);

	INSERT INTO Categories (name) SELECT 'Technology' WHERE NOT EXISTS (SELECT 1 FROM Categories WHERE name = 'Technology');
	INSERT INTO Categories (name) SELECT 'Sports' WHERE NOT EXISTS (SELECT 1 FROM Categories WHERE name = 'Sports');
	INSERT INTO Categories (name) SELECT 'Music' WHERE NOT EXISTS (SELECT 1 FROM Categories WHERE name = 'Music');
	INSERT INTO Categories (name) SELECT 'Movies' WHERE NOT EXISTS (SELECT 1 FROM Categories WHERE name = 'Movies');
	INSERT INTO Categories (name) SELECT 'Food' WHERE NOT EXISTS (SELECT 1 FROM Categories WHERE name = 'Food');
	INSERT INTO Categories (name) SELECT 'Travel' WHERE NOT EXISTS (SELECT 1 FROM Categories WHERE name = 'Travel');
	INSERT INTO Categories (name) SELECT 'Science' WHERE NOT EXISTS (SELECT 1 FROM Categories WHERE name = 'Science');
	INSERT INTO Categories (name) SELECT 'Books' WHERE NOT EXISTS (SELECT 1 FROM Categories WHERE name = 'Books');
	INSERT INTO Categories (name) SELECT 'Gaming' WHERE NOT EXISTS (SELECT 1 FROM Categories WHERE name = 'Gaming');
	INSERT INTO Categories (name) SELECT 'Health' WHERE NOT EXISTS (SELECT 1 FROM Categories WHERE name = 'Health');

	`
	
	
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	log.Println("Returned db.") 
	return db
}
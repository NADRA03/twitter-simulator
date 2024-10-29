package backend

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func Database() error {
	var err error
	database, err = sql.Open("sqlite3", "forum.db") 
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	
	table := getTable()
	_, err = database.Exec(table)
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}
	return nil
}


func getTable() string {
	return`	
	 CREATE TABLE IF NOT EXISTS users (
       	id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
        age INTEGER NOT NULL,
        gender TEXT NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS Sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);


`

}

func GetDB() *sql.DB {
	return database
}

func CloseDB() error {
	if database != nil {
		return database.Close()
	}
	return nil
}

func UserRegister(db *sql.DB,username, email, firstName, lastName, age, gender, password string) error {
	fmt.Print(username, email, firstName, lastName, age, gender, password)
	created_at:=time.Now()
	stmt, err := db.Prepare("INSERT INTO users (username, email, first_name, last_name, age, gender, password , created_at) VALUES (?, ?, ?, ?, ?, ?, ?,?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, email, firstName, lastName, age, gender, password,created_at)
	if err != nil {
		return fmt.Errorf("error executing statement: %v", err)
	}
	return nil
}
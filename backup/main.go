package main

import (
	"forum/backend"
	"log"
	"database/sql"
	"net/http"
)


var db *sql.DB

func main() {
	err := backend.Database() 
	if err != nil {
		log.Fatalf("error initializing database: %v", err)
	}
	defer backend.CloseDB()

	err = backend.DeleteExpiredSessions(backend.GetDB()) // Now this will work
	if err != nil {
		log.Printf("Error deleting expired sessions: %v", err)
	}

	
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/css"))))
    http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("frontend"))))

	http.HandleFunc("/", backend.HomeHandler)
	http.HandleFunc("/login", backend.LoginHandeler)
	http.HandleFunc("/signup", backend.SignUpHandler)
	http.HandleFunc("/home", backend.HomePageHandler)
	

	log.Println("server started on :8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
	log.Println("a")
}
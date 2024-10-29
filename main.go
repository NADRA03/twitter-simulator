package main

import (
    "twitter/backend"
    "fmt"
    "net/http"
)

// indexHandler serves index.html for any request
func indexHandler(w http.ResponseWriter, r *http.Request) {
    // Serve index.html for any route
    http.ServeFile(w, r, "./index.html")
}

func main() {
    twitter.InitializeDB() 
    // Test calling GetDB to log the message
    twitter.GetDB() // This will trigger the log
    http.HandleFunc("/", indexHandler)         
    http.HandleFunc("/log-in", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            indexHandler(w, r)
        } else if r.Method == http.MethodPost {
            twitter.LoginHandler(w, r) 
        }
    })
    http.HandleFunc("/chats", indexHandler)     
    http.HandleFunc("/chat/", indexHandler)     
    
    handler := &twitter.ChatsHandler{}
	http.HandleFunc("/chats/create", handler.CreateChatHandler)
	http.HandleFunc("/chats/addUser", handler.AddUserToChatHandler)
	http.HandleFunc("/chats/user", handler.GetUserChatsHandler)
    http.HandleFunc("/allusers", twitter.GetAllUserDetailsHandler)
    http.HandleFunc("/log-in/create-account", twitter.SignUpHandler) 
    http.HandleFunc("/User", twitter.UserDetailsHandler) 

    http.Handle("/loader.js", http.FileServer(http.Dir("."))) 
    http.Handle("/sidebar.js", http.FileServer(http.Dir("."))) 
    http.Handle("/home.js", http.FileServer(http.Dir("."))) 
    http.Handle("/sidebar2.js", http.FileServer(http.Dir("."))) 
    http.Handle("/log-in.js", http.FileServer(http.Dir(".")))
    http.Handle("/chats.js", http.FileServer(http.Dir("."))) 
    http.Handle("/chat.js", http.FileServer(http.Dir("."))) 
    http.Handle("/Profile.js", http.FileServer(http.Dir("."))) 



    http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./fonts"))))
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
    // Start the server on port 8080
    fmt.Println("Server started at :8088")
    if err := http.ListenAndServe(":8088", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
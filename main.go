package main

import (
    "twitter/backend"
    "fmt"
    "net/http"
    //"log"
)

// indexHandler serves index.html for any request
func indexHandler(w http.ResponseWriter, r *http.Request) {
    // Serve index.html for any route
    http.ServeFile(w, r, "./index.html")
    	
}

/*
func createPostHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request for %s with method %s", r.URL.Path, r.Method)
    if r.Method == http.MethodGet {
        // Serve the main page; the JS will take care of routing
        indexHandler(w, r)
    } else if r.Method == http.MethodPost {
        // Call your CreatePostHandler logic here
        twitter.CreatePostHandler(w, r)
    }
}
*/

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
    http.HandleFunc("/createPost", twitter.CreatePostHandler)
    http.HandleFunc("/addComment", twitter.AddCommentHandler)
     http.HandleFunc("/post/", twitter.PostCommentsHandler)
    http.Handle("/loader.js", http.FileServer(http.Dir("."))) 
    http.Handle("/sidebar.js", http.FileServer(http.Dir("."))) 
    http.Handle("/home.js", http.FileServer(http.Dir("."))) 
    http.Handle("/sidebar2.js", http.FileServer(http.Dir("."))) 
    http.Handle("/log-in.js", http.FileServer(http.Dir(".")))
    http.Handle("/chats.js", http.FileServer(http.Dir("."))) 
    http.Handle("/chat.js", http.FileServer(http.Dir("."))) 
    http.Handle("/Profile.js", http.FileServer(http.Dir("."))) 
    http.Handle("/post.js", http.FileServer(http.Dir("."))) 


    http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./fonts"))))
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
    // Start the server on port 8080
    fmt.Println("Server started at :8089")
    if err := http.ListenAndServe(":8089", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
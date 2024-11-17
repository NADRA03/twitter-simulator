package twitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// ChatsHandler handles requests related to chat functionalities
type ChatsHandler struct{}

// CreateChatHandler handles the creation of a new chat
func (h *ChatsHandler) CreateChatHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Chat is being created")

    // Validate the user session
    session, err := ValidateSession(w, r, db)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        log.Println("Unauthorized: Invalid session")
        return
    }

    // Parse the form data
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        log.Println("Error parsing form:", err)
        return
    }

    // Retrieve form values
    name := r.FormValue("name")   // Chat name
    bio := r.FormValue("bio")     // Chat bio
    image := r.FormValue("image") // Image URL
    chatType := r.FormValue("type") // Chat type

    log.Println("Parsed form values:", name, bio, image, chatType) // Debug log

    // Default chatType to 'private' if it's not valid
    if chatType != "private" && chatType != "group" {
        log.Println("Invalid chat type provided, defaulting to 'private'")
        chatType = "private" // Default value
    }

    var chatID int
    query := "INSERT INTO chats (chat_name, bio, type"
    args := []interface{}{name, bio, chatType} // Include chatType

    if image != "" {
        query += ", image"
        args = append(args, image)
    }
    
    query += ") VALUES (?, ?, ?"
    if image != "" {
        query += ", ?"
    }
    query += ") RETURNING id"

    // Execute the database query
    err = db.QueryRow(query, args...).Scan(&chatID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating chat: %v", err), http.StatusInternalServerError)
        log.Println("Error creating chat:", err)
        return
    }

    // Add the user to the chat
    _, err = db.Exec("INSERT INTO chat_users (chat_id, user_id, role) VALUES (?, ?, ?)", chatID, session.UserID, "admin")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error adding user to chat: %v", err), http.StatusInternalServerError)
        log.Println("Error adding user to chat:", err)
        return
    }

    // Respond with the chat ID
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{"chat_id": chatID})
}



// AddUserToChatHandler handles adding a user to a chat
func (h *ChatsHandler) AddUserToChatHandler(w http.ResponseWriter, r *http.Request) {
    session, err := ValidateSession(w, r, db)
    if err != nil {
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing form data", http.StatusBadRequest)
        return
    }

    chatID, err := strconv.Atoi(r.FormValue("chat_id"))
    if err != nil {
        http.Error(w, "Invalid chat ID", http.StatusBadRequest)
        return
    }
    
    userIDStr := r.FormValue("user_id")
    log.Println("Received user ID:", userIDStr)

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    log.Println("Parsed user ID:", userID)

    role := r.FormValue("role")
    if session.UserID != userID {
        http.Error(w, "Permission denied", http.StatusForbidden)
        return
    }

    query := "INSERT INTO chat_users (chat_id, user_id, role) VALUES (?, ?, ?)"
    _, err = db.Exec(query, chatID, userID, role)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error adding user to chat: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}




// GetUserChatsHandler handles fetching the chats for a user
func (h *ChatsHandler) GetUserChatsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Request received to send chats")

    // Validate user session
    session, err := ValidateSession(w, r, db) // Using db directly
    if err != nil {
        http.Error(w, "Session invalid or expired", http.StatusUnauthorized)
        log.Printf("Session validation error: %v", err)
        return
    }

    // Prepare the query to fetch chat details and the last message for the given user ID
    query := `
        SELECT 
            c.id AS chat_id, 
            c.chat_name, 
            c.image,
            COALESCE(m.message_text, '') AS last_message_text,
            COALESCE(m.created_at, '') AS last_message_created_at
        FROM chats c
        JOIN chat_users cu ON c.id = cu.chat_id
        LEFT JOIN (
            SELECT chat_id, message_text, created_at
            FROM messages
            WHERE (chat_id, created_at) IN (
                SELECT chat_id, MAX(created_at)
                FROM messages
                GROUP BY chat_id
            )
        ) m ON c.id = m.chat_id
        WHERE cu.user_id = ?
    `

    // Execute the query
    rows, err := db.Query(query, session.UserID)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        http.Error(w, "Error fetching user chats", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Initialize a slice to hold chats
    var chats []Chat
    for rows.Next() {
        var chat Chat
        if err := rows.Scan(&chat.ChatID, &chat.Name, &chat.ImageURL, &chat.LastMessage.MessageText, &chat.LastMessage.CreatedAt); err != nil {
            log.Printf("Error scanning chat data: %v", err)
            http.Error(w, "Error scanning chat data", http.StatusInternalServerError)
            return
        }

        // Append the chat to the list
        chats = append(chats, chat)
    }

    // Check for any errors that occurred during iteration
    if err := rows.Err(); err != nil {
        log.Printf("Error reading chat data: %v", err)
        http.Error(w, "Error reading chat data", http.StatusInternalServerError)
        return
    }

    // Set the response header to JSON and send the chats
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(chats); err != nil {
        log.Printf("Error encoding response as JSON: %v", err)
        http.Error(w, "Error encoding response as JSON", http.StatusInternalServerError)
    }

}
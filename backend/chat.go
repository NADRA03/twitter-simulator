package twitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
    "time"
    "sync"
    "database/sql"
    "github.com/gorilla/websocket" 
)

type Invitation struct {
    GroupName  string    `json:"group_name"`  
    ChatID     int       `json:"chat_id"`
    UserID     int       `json:"user_id"`
    InviterID  int       `json:"inviter_id"`
    DateTime   time.Time `json:"date_time"`
}

type ChatDetails struct {
    ChatID    int       `json:"chat_id"`
    ChatName  string    `json:"chat_name"`
    ImageURL  string    `json:"image_url"`
    Users     []User    `json:"users"`
    Messages  []Message  `json:"messages"`
}

type Message struct {
    MessageID   int             `json:"message_id"`
    UserID      int             `json:"user_id"`      // Add UserID field
    MessageText sql.NullString   `json:"message_text"`
    ImageURL    sql.NullString   `json:"image_url"`
    CreatedAt   sql.NullTime     `json:"created_at"`
}



// ChatsHandler handles requests related to chat functionalities
type ChatsHandler struct{
	connections map[*websocket.Conn]int // WebSocket connections and their associated chat IDs
	mutex       sync.Mutex
	upgrader    websocket.Upgrader
}

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
    // Validate the session
    _, err := ValidateSession(w, r, db)
    if err != nil {
        log.Println("Session validation failed:", err)
        return
    }

    // Parse the form data
    if err := r.ParseForm(); err != nil {
        log.Println("Error parsing form data:", err)
        http.Error(w, "Error parsing form data", http.StatusBadRequest)
        return
    }

    // Log all form values for inspection
    log.Println("Form values:", r.Form)

    // Parse chat_id
    chatID, err := strconv.Atoi(r.FormValue("chat_id"))
    if err != nil {
        log.Println("Invalid chat ID:", r.FormValue("chat_id"), err)
        http.Error(w, "Invalid chat ID", http.StatusBadRequest)
        return
    }

    // Parse user_id
    userIDStr := r.FormValue("user_id")
    log.Println("Received user ID string:", userIDStr)

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        log.Println("Invalid user ID:", userIDStr, err)
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    log.Println("Parsed user ID:", userID)

    // Check user permissions
    role := r.FormValue("role")
    // if session.UserID != userID {
    //     log.Println("Permission denied for user ID:", userID)
    //     http.Error(w, "Permission denied", http.StatusForbidden)
    //     return
    // }
    
    // Insert the user into chat_users
    query := "INSERT INTO chat_users (chat_id, user_id, role) VALUES (?, ?, ?)"
    _, err = db.Exec(query, chatID, userID, role)
    if err != nil {
        log.Println("Error adding user to chat:", err)
        http.Error(w, fmt.Sprintf("Error adding user to chat: %v", err), http.StatusInternalServerError)
        return
    }

    // Respond with no content on success
    log.Println("User added to chat successfully:", userID, "to chat ID:", chatID)
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


//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 

// CreateInvitationHandler handles the creation of invitations
func CreateInvitationHandler(h *ChatsHandler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
            return
        }

        var invitation Invitation
        err := json.NewDecoder(r.Body).Decode(&invitation)
        if err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Insert the invitation into the database
        query := `INSERT INTO invitations (group_name, chat_id, user_id, inviter_id, date_time) 
                  VALUES (?, ?, ?, ?, ?)`
        _, err = db.Exec(query, invitation.GroupName, invitation.ChatID, invitation.UserID, invitation.InviterID, time.Now())
        if err != nil {
            http.Error(w, "Failed to create invitation", http.StatusInternalServerError)
            return
        }

        // Respond with success
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(invitation)
    }
}

// GetChatDetailsHandler handles fetching chat details including messages and users
func (h *ChatsHandler) GetChatDetailsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Request received to fetch chat details")

    // Validate user session
    _, err := ValidateSession(w, r, db)
    if err != nil {
        http.Error(w, "Session invalid or expired", http.StatusUnauthorized)
        log.Printf("Session validation error: %v", err)
        return
    }

    // Get chat ID from query parameters
    chatIDStr := r.URL.Query().Get("chat_id")
    chatID, err := strconv.Atoi(chatIDStr)
    if err != nil || chatID <= 0 {
        http.Error(w, "Invalid chat ID", http.StatusBadRequest)
        return
    }

    // Prepare the query to fetch chat details, messages, and users
    query := `
        SELECT 
            c.id AS chat_id,
            c.chat_name,
            c.image AS image_url,
            u.id AS user_id,
            u.username AS name,
            u.image_url AS user_image_url,
            m.message_id AS message_id,
            m.user_id AS message_user_id,  -- Added user_id for messages
            m.message_text,
            m.created_at
        FROM chats c
        LEFT JOIN chat_users cu ON c.id = cu.chat_id
        LEFT JOIN users u ON cu.user_id = u.id
        LEFT JOIN messages m ON c.id = m.chat_id
        WHERE c.id = ?
    `

    // Execute the query
    rows, err := db.Query(query, chatID)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        http.Error(w, "Error fetching chat details", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Initialize chat details structure
    var chatDetails ChatDetails
    chatDetails.ChatID = chatID

    // Temporary maps to hold user details and messages
    userMap := make(map[int]User)
    var messages []Message

    for rows.Next() {
        var userID int
        var userName string
        var profileImage sql.NullString
        var messageID int
        var messageUserID int           // Variable for message user ID
        var messageText sql.NullString
        var createdAt sql.NullTime

        err := rows.Scan(&chatDetails.ChatID, &chatDetails.ChatName, &chatDetails.ImageURL, &userID, &userName, &profileImage, &messageID, &messageUserID, &messageText, &createdAt)
        if err != nil {
            log.Printf("Error scanning chat data: %v", err)
            http.Error(w, "Error scanning chat data", http.StatusInternalServerError)
            return
        }

        // Add user details to the map if not already added
        if _, exists := userMap[userID]; !exists {
            userMap[userID] = User{Id: userID, Username: userName, ImageURL: profileImage}
        }

        // Add message to messages slice if messageText is not empty
        if messageText.Valid {
            messages = append(messages, Message{
                MessageID:   messageID,
                UserID:      messageUserID,  // Include user ID for each message
                MessageText: messageText, 
                CreatedAt:   createdAt,
            })
        }
    }

    // Convert user map to slice
    var users []User
    for _, user := range userMap {
        users = append(users, user)
    }

    // Prepare response with chat details, users, and messages
    response := struct {
        ChatDetails ChatDetails `json:"chat_details"`
        Users       []User      `json:"users"`
        Messages    []Message   `json:"messages"`
    }{
        ChatDetails: chatDetails,
        Users:       users,
        Messages:    messages,
    }
    log.Printf("sent details")
    
    // Set response header and encode response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
package twitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
    "time"
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
    MessageID   int       `json:"message_id"`   
    MessageText string    `json:"message_text"`
    ImageURL    string    `json:"image_url"`    
    CreatedAt   time.Time `json:"created_at"`
}



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

func (h *ChatsHandler) AddMessageHandler(w http.ResponseWriter, r *http.Request) {
    session, err := ValidateSession(w, r, db)
    if err != nil {
        return
    }

    // Use the UserID from the session
    userID := session.UserID

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing form data", http.StatusBadRequest)
        return
    }

    chatIDStr := r.FormValue("chat_id")
    chatID, err := strconv.Atoi(chatIDStr)
    if err != nil {
        http.Error(w, "Invalid chat ID", http.StatusBadRequest)
        return
    }

    log.Println("Parsed chat ID:", chatID)

    // Check if the user is part of the chat before allowing them to send a message
    if !h.IsUserInChat(userID, chatID) {
        http.Error(w, "User is not part of this chat", http.StatusForbidden)
        return
    }

    messageText := r.FormValue("message_text")
    imageURL := r.FormValue("image_url") // Assuming this field is optional

    query := "INSERT INTO messages (chat_id, user_id, message_text, image_url) VALUES (?, ?, ?, ?)"
    _, err = db.Exec(query, chatID, userID, messageText, imageURL)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error adding message: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// Helper function to check if a user is part of the chat
func (h *ChatsHandler) IsUserInChat(userID, chatID int) bool {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM chat_users WHERE chat_id = ? AND user_id = ?)"
    err := db.QueryRow(query, chatID, userID).Scan(&exists)
    if err != nil {
        log.Println("Error checking user in chat:", err)
        return false
    }
    return exists
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
            u.user_id,
            u.username AS name,
            m.message_text,
            m.created_at
        FROM chats c
        LEFT JOIN chat_users cu ON c.id = cu.chat_id
        LEFT JOIN users u ON cu.user_id = u.user_id
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
        var messageText string
        var createdAt time.Time

        err := rows.Scan(&chatDetails.ChatID, &chatDetails.ChatName, &chatDetails.ImageURL, &userID, &userName, &messageText, &createdAt)
        if err != nil {
            log.Printf("Error scanning chat data: %v", err)
            http.Error(w, "Error scanning chat data", http.StatusInternalServerError)
            return
        }

        // Add user details to the map if not already added
        if _, exists := userMap[userID]; !exists {
            userMap[userID] = User{Id: userID, Username: userName}
        }

        // Add message to messages slice if messageText is not empty
        if messageText != "" {
            messages = append(messages, Message{MessageText: messageText, CreatedAt: createdAt})
        }
    }

    // Populate chat details users from the map
    for _, user := range userMap {
        chatDetails.Users = append(chatDetails.Users, user)
    }

    // Check for any errors that occurred during iteration
    if err := rows.Err(); err != nil {
        log.Printf("Error reading chat data: %v", err)
        http.Error(w, "Error reading chat data", http.StatusInternalServerError)
        return
    }

    // Assign messages to chat details
    chatDetails.Messages = messages

    // Set the response header to JSON and send the chat details
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(chatDetails); err != nil {
        log.Printf("Error encoding response as JSON: %v", err)
        http.Error(w, "Error encoding response as JSON", http.StatusInternalServerError)
    }
}

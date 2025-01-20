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
    ChatType  string    `json:"type"`
    ImageURL  string    `json:"image_url"`
    Users     []User    `json:"users"`
    Messages  []Message  `json:"messages"`
}

type Message struct {
    MessageID   sql.NullInt64   `json:"message_id"`   
    UserID      sql.NullInt64   `json:"user_id"`     
    MessageText sql.NullString  `json:"message_text"`
    ImageURL    sql.NullString  `json:"image_url"`
    CreatedAt   sql.NullTime    `json:"created_at"`
    SenderName  sql.NullString  `json:"sender_name"`
}


type ChatsHandler struct{
	connections map[*websocket.Conn]int 
	mutex       sync.Mutex
	upgrader    websocket.Upgrader  
    onlineUsers map[int]bool  
}


//create group chat
func (h *ChatsHandler) CreateChatHandler(w http.ResponseWriter, r *http.Request) {
    // log.Println("Chat is being created")

    session, err := ValidateSession(w, r, db)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        log.Println("Unauthorized: Invalid session")
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        log.Println("Error parsing form:", err)
        return
    }

    name := r.FormValue("name")   
    bio := r.FormValue("bio")    
    image := r.FormValue("image") 
    chatType := r.FormValue("type") 

    // log.Println("Parsed form values:", name, bio, image, chatType) 

    if chatType != "private" && chatType != "group" {
        // log.Println("Invalid chat type provided, defaulting to 'private'")
        chatType = "private" 
    }

    var chatID int
    query := "INSERT INTO chats (chat_name, bio, type"
    args := []interface{}{name, bio, chatType} 

    if image != "" {
        query += ", image"
        args = append(args, image)
    }
    
    query += ") VALUES (?, ?, ?"
    if image != "" {
        query += ", ?"
    }
    query += ") RETURNING id"

    err = db.QueryRow(query, args...).Scan(&chatID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating chat: %v", err), http.StatusInternalServerError)
        log.Println("Error creating chat:", err)
        return
    }

    _, err = db.Exec("INSERT INTO chat_users (chat_id, user_id, role) VALUES (?, ?, ?)", chatID, session.UserID, "admin")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error adding user to chat: %v", err), http.StatusInternalServerError)
        log.Println("Error adding user to chat:", err)
        return
    }

    redirectURL := fmt.Sprintf("/chat/%d", chatID)
    http.Redirect(w, r, redirectURL, http.StatusSeeOther) 
}

// private chat
func (h *ChatsHandler) CreateDirectHandler(w http.ResponseWriter, r *http.Request) {
    // log.Println("Processing chat creation or redirection")

    session, err := ValidateSession(w, r, db)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        log.Println("Unauthorized: Invalid session")
        return
    }

    otherUserIDStr := r.URL.Query().Get("id")
    otherUserID, err := strconv.Atoi(otherUserIDStr)
    if err != nil || otherUserID <= 0 {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        log.Println("Invalid user ID provided:", otherUserIDStr)
        return
    }

    var existingChatID int
    query := `
        SELECT c.id 
        FROM chats c
        JOIN chat_users cu1 ON c.id = cu1.chat_id
        JOIN chat_users cu2 ON c.id = cu2.chat_id
        WHERE c.type = 'private' AND cu1.user_id = ? AND cu2.user_id = ?`
    err = db.QueryRow(query, session.UserID, otherUserID).Scan(&existingChatID)
    if err == nil {
        // log.Println("Private chat already exists. Redirecting to chat ID:", existingChatID)
        response := map[string]string{
            "redirectUrl": fmt.Sprintf("/chat/%d", existingChatID),
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        return
    }

    // log.Println("No existing private chat found. Creating a new chat.")
    var newChatID int
    insertChatQuery := "INSERT INTO chats (chat_name, bio, type, image) VALUES (?, ?, ?, ?) RETURNING id"
    err = db.QueryRow(insertChatQuery, "Untitled Chat", "", "private", "").Scan(&newChatID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating chat: %v", err), http.StatusInternalServerError)
        log.Println("Error creating chat:", err)
        return
    }

    _, err = db.Exec(`
        INSERT INTO chat_users (chat_id, user_id, role) VALUES (?, ?, ?), (?, ?, ?)`,
        newChatID, session.UserID, "admin", 
        newChatID, otherUserID, "participant") 
    if err != nil {
        http.Error(w, fmt.Sprintf("Error adding users to chat: %v", err), http.StatusInternalServerError)
        log.Println("Error adding users to chat:", err)
        return
    }

    // log.Println("New private chat created. Redirecting to chat ID:", newChatID)
    response := map[string]string{
        "redirectUrl": fmt.Sprintf("/chat/%d", newChatID),
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


//add user to chat
func (h *ChatsHandler) AddUserToChatHandler(w http.ResponseWriter, r *http.Request) {
    _, err := ValidateSession(w, r, db)
    if err != nil {
        log.Println("Session validation failed:", err)
        return
    }

    if err := r.ParseForm(); err != nil {
        log.Println("Error parsing form data:", err)
        http.Error(w, "Error parsing form data", http.StatusBadRequest)
        return
    }

    log.Println("Form values:", r.Form)

    chatID, err := strconv.Atoi(r.FormValue("chat_id"))
    if err != nil {
        log.Println("Invalid chat ID:", r.FormValue("chat_id"), err)
        http.Error(w, "Invalid chat ID", http.StatusBadRequest)
        return
    }

    userIDStr := r.FormValue("user_id")
    log.Println("Received user ID string:", userIDStr)

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        log.Println("Invalid user ID:", userIDStr, err)
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    log.Println("Parsed user ID:", userID)

    role := r.FormValue("role")
    query := "INSERT INTO chat_users (chat_id, user_id, role) VALUES (?, ?, ?)"
    _, err = db.Exec(query, chatID, userID, role)
    if err != nil {
        log.Println("Error adding user to chat:", err)
        http.Error(w, fmt.Sprintf("Error adding user to chat: %v", err), http.StatusInternalServerError)
        return
    }

    // log.Println("User added to chat successfully:", userID, "to chat ID:", chatID)
    w.WriteHeader(http.StatusNoContent)
}

// All chats
func (h *ChatsHandler) GetUserChatsHandler(w http.ResponseWriter, r *http.Request) {
    // log.Println("Request received to send chats")

    session, err := ValidateSession(w, r, db)
    if err != nil {
        http.Error(w, "Session invalid or expired", http.StatusUnauthorized)
        log.Printf("Session validation error: %v", err)
        return
    }

    query := `
        SELECT 
            c.id AS chat_id, 
            CASE 
                WHEN c.type = 'private' THEN (
                    SELECT u.username 
                    FROM users u
                    JOIN chat_users cu ON u.id = cu.user_id
                    WHERE cu.chat_id = c.id AND u.id != ?
                )
                ELSE c.chat_name
            END AS chat_name,
            CASE 
                WHEN c.type = 'private' THEN (
                    SELECT u.image_url 
                    FROM users u
                    JOIN chat_users cu ON u.id = cu.user_id
                    WHERE cu.chat_id = c.id AND u.id != ?
                )
                ELSE c.image
            END AS chat_image_url,
            CASE 
                WHEN c.type = 'private' THEN (
                    SELECT u.status 
                    FROM users u
                    JOIN chat_users cu ON u.id = cu.user_id
                    WHERE cu.chat_id = c.id AND u.id != ?
                )
                ELSE NULL
            END AS chat_status,
            c.type AS chat_type,
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

    rows, err := db.Query(query, session.UserID, session.UserID, session.UserID, session.UserID)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        http.Error(w, "Error fetching user chats", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var chats []Chat
    for rows.Next() {
        var chat Chat
        var chatImageURL sql.NullString
        var chatStatus sql.NullString
        if err := rows.Scan(&chat.ChatID, &chat.Name, &chatImageURL, &chatStatus, &chat.Type, &chat.LastMessage.MessageText, &chat.LastMessage.CreatedAt); err != nil {
            log.Printf("Error scanning chat data: %v", err)
            http.Error(w, "Error scanning chat data", http.StatusInternalServerError)
            return
        }

        if chatImageURL.Valid {
            chat.ImageURL = chatImageURL.String
        } else {
            chat.ImageURL = ""
        }

        if chatStatus.Valid {
            chat.Status = chatStatus.String
        } else {
            chat.Status = ""
        }

        chats = append(chats, chat)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error reading chat data: %v", err)
        http.Error(w, "Error reading chat data", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(chats); err != nil {
        log.Printf("Error encoding response as JSON: %v", err)
        http.Error(w, "Error encoding response as JSON", http.StatusInternalServerError)
    }
}



//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 

//invite people
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

        query := `INSERT INTO invitations (group_name, chat_id, user_id, inviter_id, date_time) 
                  VALUES (?, ?, ?, ?, ?)`
        _, err = db.Exec(query, invitation.GroupName, invitation.ChatID, invitation.UserID, invitation.InviterID, time.Now())
        if err != nil {
            http.Error(w, "Failed to create invitation", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(invitation)
    }
}


//all chat details and first 10 messages
func (h *ChatsHandler) GetChatDetailsHandler(w http.ResponseWriter, r *http.Request) {
    // log.Println("Request received to fetch chat details")

    _, err := ValidateSession(w, r, db)
    if err != nil {
        http.Error(w, "Session invalid or expired", http.StatusUnauthorized)
        log.Printf("Session validation error: %v", err)
        return
    }

    chatIDStr := r.URL.Query().Get("chat_id")
    chatID, err := strconv.Atoi(chatIDStr)
    if err != nil || chatID <= 0 {
        http.Error(w, "Invalid chat ID", http.StatusBadRequest)
        return
    }

    query := `
        SELECT 
            c.id AS chat_id,
            c.chat_name,
            c.type,
            c.image AS image_url,
            u.id AS user_id,
            u.username AS name,
            u.image_url AS user_image_url
        FROM chats c
        LEFT JOIN chat_users cu ON c.id = cu.chat_id
        LEFT JOIN users u ON cu.user_id = u.id
        WHERE c.id = ?`

    rows, err := db.Query(query, chatID)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        http.Error(w, "Error fetching chat details", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var chatDetails ChatDetails
    chatDetails.ChatID = chatID

    userMap := make(map[int]User)

    for rows.Next() {
        var userID int
        var userName string
        var profileImage sql.NullString

        err := rows.Scan(&chatDetails.ChatID, &chatDetails.ChatName, &chatDetails.ChatType,&chatDetails.ImageURL, &userID, &userName, &profileImage)
        if err != nil {
            log.Printf("Error scanning chat data: %v", err)
            http.Error(w, "Error scanning chat data", http.StatusInternalServerError)
            return
        }

        if _, exists := userMap[userID]; !exists {
            userMap[userID] = User{Id: userID, Username: userName, ImageURL: profileImage}
        }
    }

    messagesQuery := `
        SELECT 
            m.message_id,
            m.user_id,
            m.message_text,
            m.created_at
        FROM messages m
        WHERE m.chat_id = ?
        ORDER BY m.created_at DESC
        LIMIT 10`

    messageRows, err := db.Query(messagesQuery, chatID)
    if err != nil {
        log.Printf("Error fetching messages: %v", err)
        http.Error(w, "Error fetching messages", http.StatusInternalServerError)
        return
    }
    defer messageRows.Close()

    var messages []Message
    for messageRows.Next() {
        var message Message
        err := messageRows.Scan(&message.MessageID, &message.UserID, &message.MessageText, &message.CreatedAt)
        if err != nil {
            log.Printf("Error scanning message data: %v", err)
            http.Error(w, "Error scanning message data", http.StatusInternalServerError)
            return
        }
        messages = append(messages, message)
    }

    for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
        messages[i], messages[j] = messages[j], messages[i]
    }

    var users []User
    for _, user := range userMap {
        users = append(users, user)
    }

    response := struct {
        ChatDetails ChatDetails `json:"chat_details"`
        Users       []User      `json:"users"`
        Messages    []Message   `json:"messages"`
    }{
        ChatDetails: chatDetails,
        Users:       users,
        Messages:    messages,
    }

    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Error sending response", http.StatusInternalServerError)
    } else {
        // log.Println("Sent chat details successfully")
    }
}


//more 10 messages
func (h *ChatsHandler) GetMoreMessagesHandler(w http.ResponseWriter, r *http.Request) {
    // log.Println("Request received to fetch more messages")

    _, err := ValidateSession(w, r, db)
    if err != nil {
        http.Error(w, "Session invalid or expired", http.StatusUnauthorized)
        log.Printf("Session validation error: %v", err)
        return
    }

    chatIDStr := r.URL.Query().Get("chatId")
    chatID, err := strconv.Atoi(chatIDStr)
    if err != nil || chatID <= 0 {
        http.Error(w, "Invalid chat ID", http.StatusBadRequest)
        return
    }

    lastMessageIDStr := r.URL.Query().Get("lastMessageId")
    lastMessageID, err := strconv.Atoi(lastMessageIDStr)
    if err != nil || lastMessageID <= 0 {
        http.Error(w, "Invalid last message ID", http.StatusBadRequest)
        return
    }

    query := `
        SELECT 
            m.message_id,
            m.user_id,
            m.message_text,
            m.created_at
        FROM messages m
        WHERE m.chat_id = ? AND m.message_id < ?
        ORDER BY m.created_at DESC
        LIMIT 10`

    rows, err := db.Query(query, chatID, lastMessageID)
    if err != nil {
        log.Printf("Error fetching messages: %v", err)
        http.Error(w, "Error fetching messages", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var messages []Message
    for rows.Next() {
        var message Message
        err := rows.Scan(&message.MessageID, &message.UserID, &message.MessageText, &message.CreatedAt)
        if err != nil {
            log.Printf("Error scanning message data: %v", err)
            http.Error(w, "Error scanning message data", http.StatusInternalServerError)
            return
        }
        messages = append(messages, message)
    }
    
    for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
        messages[i], messages[j] = messages[j], messages[i]
    }


    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(struct {
        Messages []Message `json:"messages"`
    }{Messages: messages})
    if err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Error sending response", http.StatusInternalServerError)
    } else {
        // log.Println("Sent more messages successfully")
    }
}

func (h *ChatsHandler) GetAllLastMessagesHandler(w http.ResponseWriter, r *http.Request) {
    // log.Println("Request received to fetch last messages of all chats")

    // Validate session and get the session details
    session, err := ValidateSession(w, r, db)
    if err != nil {
        log.Printf("Session validation error: %v", err)
        http.Error(w, "Session invalid or expired", http.StatusUnauthorized)
        return
    }

    // Get user ID from the session
    userID := session.UserID

    // Query to get the last message for each chat the user is part of, excluding messages from the current user
    query := `
        SELECT 
            m.message_text,
            u.username AS sender_name
        FROM chats c
        LEFT JOIN chat_users cu ON c.id = cu.chat_id
        LEFT JOIN messages m ON c.id = m.chat_id
        LEFT JOIN users u ON m.user_id = u.id
        WHERE cu.user_id = ? 
          AND m.created_at >= datetime('now','-6 seconds') 
          AND m.user_id != ?  -- Exclude messages written by the current user
        ORDER BY m.created_at DESC`

    rows, err := db.Query(query, userID, userID)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        http.Error(w, "Error fetching messages", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var lastMessages []struct {
        MessageText string `json:"message_text"`
        SenderName  string `json:"sender_name"`
    }

    // Loop through rows to build the response
    for rows.Next() {
        var message struct {
            MessageText string `json:"message_text"`
            SenderName  string `json:"sender_name"`
        }

        err := rows.Scan(&message.MessageText, &message.SenderName)
        if err != nil {
            log.Printf("Error scanning message data: %v", err)
            http.Error(w, "Error scanning message data", http.StatusInternalServerError)
            return
        }

        lastMessages = append(lastMessages, message)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Error iterating over rows: %v", err)
        http.Error(w, "Error processing messages", http.StatusInternalServerError)
        return
    }

    // Send the response
    w.Header().Set("Content-Type", "application/json")
    if len(lastMessages) == 0 {
        err = json.NewEncoder(w).Encode(struct {
            Message string `json:"message"`
        }{Message: "No new messages in the last 3 seconds"})
    } else {
        err = json.NewEncoder(w).Encode(struct {
            LastMessages []struct {
                MessageText string `json:"message_text"`
                SenderName  string `json:"sender_name"`
            } `json:"last_messages"`
        }{LastMessages: lastMessages})
    }

    if err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Error sending response", http.StatusInternalServerError)
    } else {
        // log.Println("Sent last messages successfully")
        // log.Println(lastMessages)
    }
}


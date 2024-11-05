package twitter

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/websocket"
	"sync"
)

func NewChatsHandler() *ChatsHandler {
	return &ChatsHandler{
		connections: make(map[*websocket.Conn]int), // Ensure this is initialized
		mutex:      sync.Mutex{}, // Initialize the mutex
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow requests from specific origins
				origin := r.Header.Get("Origin")
				allowedOrigins := []string{
					"http://localhost:3000", 
					"http://localhost:8088",
				}

				for _, allowed := range allowedOrigins {
					if allowed == origin {
						fmt.Println("new chat handler completed")
						return true
					}
				}
				return false
			},
		},
	}
}

// AddMessageHandler adds a message to the database and broadcasts it
func (h *ChatsHandler) AddMessageToChat(chatID int, messageText string, imageURL string, userID int) {
	// Verify if the user is part of the chat
	if !h.IsUserInChat(userID, chatID) {
		log.Println("User is not part of this chat, message not saved.")
		return
	}

	// Insert message into the database
	query := "INSERT INTO messages (chat_id, user_id, message_text, image_url) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, chatID, userID, messageText, imageURL)
	if err != nil {
		logError("AddMessageToChat - Error adding message to database", err)
		return
	}

	// Confirm that a row was inserted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logError("AddMessageToChat - Error retrieving rows affected", err)
		return
	}
	if rowsAffected == 0 {
		log.Println("No rows were inserted, message was not saved to the database.")
		return
	}

	log.Printf("Message successfully saved to database with chatID %d, userID %d", chatID, userID)

	// Optionally broadcast the new message to other clients
	message := struct {
		UserID      int    `json:"user_id"`
		ChatID      int    `json:"chat_id"`
		MessageText string `json:"message_text"`
		ImageURL    string `json:"image_url"`
	}{
		UserID:      userID,
		ChatID:      chatID,
		MessageText: messageText,
		ImageURL:    imageURL,
	}

	h.BroadcastMessageToChat(chatID, message)
}



// WebSocketHandler to handle WebSocket connections
func (h *ChatsHandler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Attempting to connect to WebSocket.")

	// Step 1: Get chatID from URL and parse it
	chatIDStr := r.URL.Query().Get("chat_id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		logError("WebSocketHandler - strconv.Atoi: Invalid chat ID format", err)
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	fmt.Printf("Parsed chatID: %d\n", chatID)

	// Step 2: Log the origin and attempt to validate the session
	log.Println("WebSocket connection attempt from origin:", r.Header.Get("Origin"))

	// Step 3: Validate session and retrieve userID
	session, err := ValidateSession(w, r, db)
	if err != nil {
		logError("WebSocketHandler - ValidateSession: Session validation failed", err)
		http.Error(w, "Invalid session", http.StatusForbidden)
		return
	}
	userID := session.UserID
	fmt.Printf("Session validated, userID: %d\n", userID)

	// Step 4: Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logError("WebSocketHandler - upgrader.Upgrade: WebSocket upgrade failed", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	fmt.Println("WebSocket connection successfully upgraded.")

	// Step 5: Register connection in handler and start listening for messages
	h.mutex.Lock()
	h.connections[conn] = chatID
	h.mutex.Unlock()
	fmt.Println("Connection registered with chat ID:", chatID)

	// Start listenForMessages in a separate goroutine
	fmt.Println("Starting listenForMessages goroutine.")
	go func() {
		err := h.listenForMessages(conn, userID)
		if err != nil {
			fmt.Printf("listenForMessages error for userID %d: %v\n", userID, err)
		}
	}()
}


// BroadcastMessageToChat broadcasts a message to all users in a specified chat
func (h *ChatsHandler) BroadcastMessageToChat(chatID int, message interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for conn, userChatID := range h.connections {
		if userChatID == chatID {
			if err := conn.WriteJSON(message); err != nil {
				log.Println("Error broadcasting:", err)
				conn.Close()
				delete(h.connections, conn)
			}
		}
	}
}

// listenForMessages listens for WebSocket messages or detects disconnects
func (h *ChatsHandler) listenForMessages(conn *websocket.Conn, userID int) error {
	fmt.Println("Entering listenForMessages for userID:", userID)
	defer func() {
		fmt.Println("Closing connection in listenForMessages for userID:", userID)
		h.mutex.Lock()
		delete(h.connections, conn)
		h.mutex.Unlock()
		conn.Close()
	}()

	for {
		var messageData struct {
			Action string `json:"action"`
			Data   struct {
				ChatID      string `json:"chatId"` // Change to string
				MessageText string `json:"messageText"`
				ImageURL    string `json:"imageUrl"`
			} `json:"data"`
		}

		// Read JSON message from WebSocket
		err := conn.ReadJSON(&messageData)
		if err != nil {
			log.Printf("listenForMessages - Error reading JSON for userID %d: %v\n", userID, err)
			return fmt.Errorf("error reading JSON: %w", err)
		}

		log.Printf("Received message data from userID %d: %+v\n", userID, messageData)

		// Convert ChatID from string to int
		chatID, err := strconv.Atoi(messageData.Data.ChatID)
		if err != nil {
			log.Printf("Invalid chatId format for userID %d: %v\n", userID, err)
			continue // Skip processing this message if chatId is invalid
		}

		// Check if action is "add_message" and process it
		if messageData.Action == "add_message" {
			fmt.Println("Action is add_message, calling AddMessageToChat")
			h.AddMessageToChat(chatID, messageData.Data.MessageText, messageData.Data.ImageURL, userID)
		}
	}
}


// IsUserInChat checks if a user is part of a chat
func (h *ChatsHandler) IsUserInChat(userID, chatID int) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM chat_users WHERE chat_id = ? AND user_id = ?)"
	err := db.QueryRow(query, chatID, userID).Scan(&exists)
	if err != nil {
		logError("IsUserInChat - QueryRow", err)
		return false
	}
	return exists
}

// logError logs the error with a function name for context
func logError(funcName string, err error) {
	log.Printf("Error in %s: %v\n", funcName, err)
}
package twitter

import(
	"database/sql"
	"fmt"
)

// Define the struct for LastMessage
type LastMessage struct {
    MessageText string `json:"message_text"`
    CreatedAt   string `json:"created_at"`
}

// Define the struct for Chat
type Chat struct {
    ChatID      int       `json:"chat_id"`
    Name        string    `json:"name"`
    ImageURL    string    `json:"image"`
    LastMessage LastMessage `json:"last_message"`
}

func CreateChat(db *sql.DB, bio string, image string, chatType string) (int, error) {
    var chatID int
    var query string

    // Build the query based on whether the image is provided
    if image != "" {
        query = "INSERT INTO chats (bio, image_url, type) VALUES (?, ?, ?) RETURNING id"
    } else {
        query = "INSERT INTO chats (bio, type) VALUES (?, ?) RETURNING id"
    }

    // Execute the query and scan the chat ID
    err := db.QueryRow(query, bio, image, chatType).Scan(&chatID)
    if err != nil {
        return 0, fmt.Errorf("error creating chat: %v", err)
    }
    return chatID, nil
}

func AddUserToChat(db *sql.DB, chatID int, userID int, role string) error {
    query := "INSERT INTO chat_users (chat_id, user_id, role) VALUES (?, ?, ?)"
    _, err := db.Exec(query, chatID, userID, role)
    if err != nil {
        return fmt.Errorf("error adding user to chat: %v", err)
    }
    return nil
}

func GetUserChats(db *sql.DB, userID int) ([]Chat, error) {
    // Prepare the query to fetch chat details and the last message for the given user ID
    query := `
        SELECT 
            c.id AS chat_id, 
            c.chat_name, 
            c.image_url,
            m.message_text AS last_message_text,
            m.created_at AS last_message_created_at
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

    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("error fetching user chats: %v", err)
    }
    defer rows.Close()

    var chats []Chat

    for rows.Next() {
        var chat Chat
        if err := rows.Scan(&chat.ChatID, &chat.Name, &chat.ImageURL, &chat.LastMessage.MessageText, &chat.LastMessage.CreatedAt); err != nil {
            return nil, fmt.Errorf("error scanning chat row: %v", err)
        }
        chats = append(chats, chat)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over chat rows: %v", err)
    }

    return chats, nil
}





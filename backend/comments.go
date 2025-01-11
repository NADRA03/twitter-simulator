package twitter

import (
	"time"
    "database/sql"
	 "fmt"
    "net/http"	
	"encoding/json"
    "strconv"
    "bytes"
    "io"
)

type Comments struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    ImageURL  sql.NullString  `json:"image_url"` 
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    ComUsername  string  `json:"username"`
}

type CountReactionC struct {
    CommentcounterID    int
    Likecount1 int       
    Dislikecount1 int    
}

func GetCommentsByPostID(postID int) ([]Comments, error) {
    // Update query to join users table and fetch image_url
    query := `
        SELECT c.id, c.post_id, c.user_id, c.content, c.username, c.created_at, u.image_url
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
    `

    rows, err := db.Query(query, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []Comments
    for rows.Next() {
        var comment Comments
        err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.ComUsername, &comment.CreatedAt, &comment.ImageURL)
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }
    return comments, nil
}


func GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
    // Extract post ID from the query string
    postID := r.URL.Query().Get("post_id")
    if postID == "" {
        http.Error(w, "Missing post ID", http.StatusBadRequest)
        return
    }

    pid, err := strconv.Atoi(postID)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    // Get comments from the database
    comments, err := GetCommentsByPostID(pid)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error fetching comments: %v", err), http.StatusInternalServerError)
        return
    }

    // Return JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    if len(comments) == 0 {
        json.NewEncoder(w).Encode(map[string]string{"message": "No comments available"})
        return
    }

    if err := json.NewEncoder(w).Encode(comments); err != nil {
        http.Error(w, fmt.Sprintf("Error encoding comments to JSON: %v", err), http.StatusInternalServerError)
        return
    }
}

func CreateComment(userID, postID int, content string , username string) error {
    createdAt := time.Now().Format(time.RFC3339)
    fmt.Println(username)
    stmt, err := db.Prepare("INSERT INTO Comments (user_id, post_id, content,  username, created_at ) VALUES (?, ?, ?, ?, ?)")
    if err != nil {
        return fmt.Errorf("error preparing insert statement: %v", err)
    }
    defer stmt.Close()
    _, err = stmt.Exec(userID, postID, content, username,createdAt)
    if err != nil {
        return fmt.Errorf("error executing insert statement: %v", err)
    }
    return nil
}

func CreateComm(w http.ResponseWriter, r *http.Request) {
    // Log the raw request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }
    fmt.Println("Received body:", string(body))

    // Reset the request body for decoding
    r.Body = io.NopCloser(bytes.NewReader(body))

    var commentData struct {
        UserID   int    `json:"user_id"`
        PostID   int    `json:"post_id"`
        Content  string `json:"content"`
        Username string `json:"username"`
    }

    // Decode JSON into struct
    if err := json.NewDecoder(r.Body).Decode(&commentData); err != nil {
        fmt.Printf("Decode error: %v\n", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    fmt.Printf("Parsed comment data: %+v\n", commentData)

    // Create the comment
    err = CreateComment(commentData.UserID, commentData.PostID, commentData.Content, commentData.Username)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating comment: %v", err), http.StatusInternalServerError)
        return
    }

    // Return success response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Comment created successfully"})
}


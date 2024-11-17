package twitter

import (
	"time"
	"database/sql"
	 "fmt"
)

type Comments struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    ComUsername  string  `json:"username"`
}

type CountReactionC struct {
    CommentcounterID    int
    Likecount1 int       
    Dislikecount1 int    
}

func GetPostByID(db *sql.DB, postID int) (Posts, error) {
	var post Posts
    err := db.QueryRow("SELECT id, title, content, user_id, created_at FROM posts WHERE id = ?", postID).
        Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt)
    if err != nil {
        return post, err
    }
    return post, nil
}

func GetCommentsByPostID(db *sql.DB, postID int) ([]Comments, error) {
    rows, err := db.Query("SELECT id, post_id, user_id, content,username ,created_at  FROM comments WHERE post_id = ?", postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var comments []Comments
    for rows.Next() {
        var comment Comments
        err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content,&comment.ComUsername, &comment.CreatedAt,)
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }
    return comments, nil
}

func CreateComment(db *sql.DB, userID, postID int, content string , username string) error {
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

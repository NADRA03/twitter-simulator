package twitter


import(
	"database/sql"
	"time"
	"fmt"
	"strconv"	
)

type Posts struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user-id"`
	Username      string
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created-at"`
	CategoriesNames []string
}

type UserReaction struct {
	UserID    int
	PostID    int 
	Reaction string
}

type MainPostCounter struct {
	PostID               int
	MainPostLikecount    int
	MainPostDislikecount int
}

func CreatePost(db *sql.DB, userID int, title, content string ) error {
	createdAt := time.Now().Format(time.RFC3339)
	stmt, err := db.Prepare("INSERT INTO Posts (user_id, title, content, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(userID, title, content, createdAt)
	if err != nil {
		return fmt.Errorf("error executing insert statement: %v", err)
	}
	return nil
}

func GetAllPosts(db *sql.DB) ([]Posts, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM Posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`
	rows, err := db.Query(query) 
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %v", err)
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func GetMyPosts(db *sql.DB, userID int) ([]Posts, error) {
	rows, err := db.Query(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at
		FROM Posts p
		JOIN Users u ON p.user_id = u.id
		Where user_id = ?
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %v", err)
	}
	defer rows.Close()
	var posts []Posts
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostID(db *sql.DB, userID int, title string, content string, createdAt string) (int, error) {
	DB := db 
	stmt, err := DB.Prepare("SELECT id FROM Posts WHERE user_id = ? AND title = ? AND content = ? AND created_at = ?")
	if err != nil {
		return 0, fmt.Errorf("error preparing select statement: %v", err)
	}
	defer stmt.Close()
	var postID int
	err = stmt.QueryRow(userID,title,content,createdAt).Scan(&postID)
	if err != nil {
		return 0, fmt.Errorf("error querying post ID: %v", err)
	}
	return postID, nil
}

// GetPostsByCategories retrieves posts that match any of the given category IDs
func GetPostsByCategories(db *sql.DB, categoryIDs []int) ([]Posts, error) {
	if len(categoryIDs) == 0 {
		return GetAllPosts(db)
	}

	query := `
    SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at
    FROM Posts p
    JOIN Users u ON p.user_id = u.id
    JOIN PostCategories pc ON p.id = pc.post_id
    WHERE pc.category_id IN (
	`
	// Append placeholders for the number of categories
	for i :=0 ; i < len(categoryIDs); i++ {
		if i > 0 {
			query += ","
		}
		query += "?"
	}

	query += " )GROUP BY p.id ORDER BY p.created_at DESC"

	rows, err := db.Query(query, convertToInterfaceSlice(categoryIDs)...)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %v", err)
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}
		posts = append(posts, post)
	}
	return posts, nil

}

func GetMyPostsByCategories(db *sql.DB, categoryIDs []int, userID int) ([]Posts, error) {
	if len(categoryIDs) == 0 {
		return GetAllPosts(db)
	}
	// convert int into string
	strUserId := strconv.Itoa(userID)

	// Prepare a query with placeholders for the category IDs
	query := `
        SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at
        FROM Posts p
        JOIN Users u ON p.user_id = u.id
        JOIN PostCategories pc ON p.id = pc.post_id
        WHERE pc.category_id IN (`

	// Append placeholders for the number of categories
	for i := 0; i < len(categoryIDs); i++ {
		if i > 0 {
			query += ", "
		}
		query += "?"
	}
	query += " ) AND user_id = " + strUserId + " GROUP BY p.id ORDER BY p.created_at DESC"

	rows, err := db.Query(query, convertToInterfaceSlice(categoryIDs)...)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %v", err)
	}
	defer rows.Close()

	var posts []Posts
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// Helper function to convert int slice to interface slice
func convertToInterfaceSlice(ints []int) []interface{} {
	iface := make([]interface{}, len(ints))
	for i, v := range ints {
		iface[i] = v
	}
	return iface
}

func RemoveDuplicatePosts(posts []Posts) []Posts {
	seen := make(map[int]bool)
	uniq := []Posts{}
	for _, post := range posts {
		if !seen[post.ID] {
			uniq = append(uniq, post)
			seen[post.ID] = true
		}
	}
	return uniq
}

func CreatePostCategory(db *sql.DB, postID int, categoryID int) error {
	stmt, err := db.Prepare("INSERT INTO PostCategories (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(postID, categoryID)
	if err != nil {
		return fmt.Errorf("error executing insert statement: %v", err)
	}
	return nil
}




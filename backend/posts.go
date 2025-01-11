package twitter


import(
	"database/sql"
	"time"
	"fmt"
	"strconv"
	"net/http"	
	"encoding/json"
	"log"
)

type Posts struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user-id"`
	ImageURL  sql.NullString  `json:"image_url"` 
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

func CreatePost(w http.ResponseWriter, r *http.Request) {
    // Ensure it's a POST request
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Parse JSON from the request body
    var post struct {
        UserID  int    `json:"userID"`
        Title   string `json:"title"`
        Content string `json:"content"`
    }

    // Decode the JSON body
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&post); err != nil {
        log.Printf("Error decoding request body: %v", err) // Log the error
        http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Validate inputs
    if post.Title == "" || post.Content == "" {
        http.Error(w, "Title and content cannot be empty", http.StatusBadRequest)
        return
    }

    // Prepare the SQL query to insert the post
    stmt, err := db.Prepare("INSERT INTO Posts (user_id, title, content, created_at) VALUES (?, ?, ?, ?)")
    if err != nil {
        log.Printf("Error preparing insert statement: %v", err) // Log the error
        http.Error(w, fmt.Sprintf("Error preparing insert statement: %v", err), http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    // Execute the insert statement
    _, err = stmt.Exec(post.UserID, post.Title, post.Content, time.Now().Format(time.RFC3339))
    if err != nil {
        log.Printf("Error executing insert statement: %v", err) // Log the error
        http.Error(w, fmt.Sprintf("Error executing insert statement: %v", err), http.StatusInternalServerError)
        return
    }

    // Return a success message
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Post created successfully"))
}


func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	// Print to confirm the function is being called
	fmt.Println("Fetching all posts")

	// SQL query to get posts, including the image_url field
	query := `
		SELECT p.id, p.user_id, u.username, u.image_url, p.title, p.content, p.created_at
		FROM Posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying posts: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the posts
	var posts []Posts

	// Loop through the rows and scan data into the posts slice
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.ImageURL, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning post: %v", err), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	// Return the response
	if len(posts) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "No posts available"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding posts to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

func GetMyPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching posts for a specific user")

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// SQL query to get posts, including the image_url field
	query := `
		SELECT p.id, p.user_id, u.username, u.image_url, p.title, p.content, p.created_at
		FROM Posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
	`

	// Execute the query
	rows, err := db.Query(query, uid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying posts: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the posts
	var posts []Posts

	// Loop through the rows and scan data into the posts slice
	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.ImageURL, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning post: %v", err), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	// Return the response
	if len(posts) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "No posts available"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding posts to JSON: %v", err), http.StatusInternalServerError)
		return
	}
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

func GetPostById(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Fetching post details for a specific post ID")

    // Get post ID from URL query parameters
    postID := r.URL.Query().Get("id")
    if postID == "" {
        http.Error(w, "Missing post ID", http.StatusBadRequest)
        return
    }

    // Convert postID from string to integer
    pid, err := strconv.Atoi(postID)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    // SQL query to get a specific post by its ID, including the image_url field
    query := `
        SELECT p.id, p.user_id, u.username, u.image_url, p.title, p.content, p.created_at
        FROM Posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.id = ?
    `

    // Execute the query to get the post by ID
    row := db.QueryRow(query, pid)

    // Create a variable to store the post
    var post Posts

    // Scan the data into the post variable
    if err := row.Scan(&post.ID, &post.UserID, &post.Username, &post.ImageURL, &post.Title, &post.Content, &post.CreatedAt); err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Post not found", http.StatusNotFound)
        } else {
            http.Error(w, fmt.Sprintf("Error scanning post: %v", err), http.StatusInternalServerError)
        }
        return
    }

    // Return the response with the post details
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(post); err != nil {
        http.Error(w, fmt.Sprintf("Error encoding post to JSON: %v", err), http.StatusInternalServerError)
        return
    }
}

// GetPostsByCategories retrieves posts that match any of the given category IDs
func GetPostsByCategories(db *sql.DB, categoryIDs []int) ([]Posts, error) {
	// if len(categoryIDs) == 0 {
	// 	return GetAllPosts()
	// }

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
	// if len(categoryIDs) == 0 {
	// 	return GetAllPosts()
	// }
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
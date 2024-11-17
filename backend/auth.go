package twitter

import(
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
    "strings"
    //"time"
    "strconv"
	"text/template"
	
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/log-in" {
		http.NotFound(w,r)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		formType := r.FormValue("form")

		if formType == "login" {
			username := r.FormValue("loginusername")
			password := r.FormValue("loginpassword")
			fmt.Printf("Received login: Username: %s, Password: %s\n", username, password)
			exiest, err := LoginUser(db, username, password)
			if err != nil {
				log.Printf("Error logging user: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !exiest {
				log.Println("Invalid username or password")
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			}

			if exiest {
				log.Println("you are logged in")
				userID, _ := GetUserID(db, username)
				err = CreateSession(w,r,db,userID)
				if err != nil {
					log.Printf("Error creating session: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

		}

	}

}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/log-in/create-account" {
        http.NotFound(w, r)
        return
    }

    if r.Method == http.MethodPost {
        err := r.ParseForm()
        if err != nil {
            http.Error(w, "Failed to parse form", http.StatusBadRequest)
            return
        }
        
        db := GetDB()
        formType := r.FormValue("form")
        
        if formType == "signup" {
            username := r.FormValue("username")
            password := r.FormValue("password")
            email := r.FormValue("email")
            firstname := r.FormValue("firstname")
            lastname := r.FormValue("lastname")
            age := r.FormValue("age")
            gender := r.FormValue("gender")

            // Validation logic
            if err := ValidUsername(db, username); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            if err := ValidPassword(password); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            if email == "" {
                http.Error(w, "Email cannot be empty.", http.StatusBadRequest)
                return
            }

            // Check if username or email already exists
            exists, err := UserExists(db, username, email)
            if err != nil {
                log.Printf("Error checking user existence: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            if exists {
                http.Error(w, "Username or email already exists.", http.StatusConflict)
                return
            }

            hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            err = UserRegister(db, username, email, firstname, lastname, age, gender, string(hashedPassword))
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            userID, _ := GetUserID(db, username)
            err = CreateSession(w, r, db, userID)
            if err != nil {
                log.Printf("Error creating session: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }
    }
}


type User struct {
	Username string `json:"username"`
	ImageURL sql.NullString `json:"image_url"`
	Id       int    `json:"id"`
    Email    string `json:"email"`      
}

// GetAllUserDetailsHandler retrieves all user details and sends them as a JSON response
func GetAllUserDetailsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Loading users")

    // Define the query to fetch all user details including ID
    query := `SELECT username, image_url, id FROM users`
    rows, err := db.Query(query)
    if err != nil {
        log.Printf("Error querying all users: %v\n", err) // Log the specific error
        http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []User // Use a slice to hold User structs
    for rows.Next() {
        var username string
        var imageUrl sql.NullString
        var id int // Variable to hold user ID

        // Scan the row into the variables
        if err := rows.Scan(&username, &imageUrl, &id); err != nil {
            log.Printf("Error scanning user row: %v\n", err) // Log the specific error
            http.Error(w, "Failed to scan user details", http.StatusInternalServerError)
            return
        }

        // Create a User struct and populate it
        user := User{
            Username: username,
            ImageURL: imageUrl,
            Id: id, // Add ID to the User struct
        }
        if !user.ImageURL.Valid {
            user.ImageURL.String = "" // Set to empty string if NULL
        }

        users = append(users, user) // Append user to the slice
    }

    // Check for errors encountered during iteration
    if err := rows.Err(); err != nil {
        log.Printf("Error during rows iteration: %v\n", err) // Log the specific error
        http.Error(w, "Failed during row iteration", http.StatusInternalServerError)
        return
    }

    // Set the response header to JSON
    w.Header().Set("Content-Type", "application/json")

    // Send the users as a JSON response
    if err := json.NewEncoder(w).Encode(users); err != nil {
        log.Printf("Error encoding JSON response: %v\n", err) // Log the specific error
        http.Error(w, "Failed to encode response", http.StatusInternalServerError) // Handle JSON encoding errors
        return
    }

    log.Println("Users sent successfully")
}


func UserDetailsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("UserDetailsHandler: Profile is loading")

    // Get the session
    session, err := GetSession(r, db)
    if err != nil {
        log.Printf("UserDetailsHandler: Error getting session: %v", err)
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    log.Printf("UserDetailsHandler: Session retrieved for UserID %d", session.UserID)

    // Query the user details based on session's UserID
    var user User
    err = db.QueryRow("SELECT id, username, email, image_url FROM Users WHERE id = ?", session.UserID).
        Scan(&user.Id, &user.Username, &user.Email, &user.ImageURL)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("UserDetailsHandler: User not found in database")
            http.Error(w, "user not found", http.StatusNotFound)
            return
        }
        log.Printf("UserDetailsHandler: Error querying database: %v", err)
        http.Error(w, "error querying database", http.StatusInternalServerError)
        return
    }

    // Check if ImageURL is valid and set it to an empty string if it is NULL
    if !user.ImageURL.Valid {
        user.ImageURL.String = "" // Set to empty string if NULL
    }

    log.Printf("UserDetailsHandler: User details retrieved: %+v", user)

    // Send the user details as JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user); err != nil {
        log.Printf("UserDetailsHandler: Error encoding user data to JSON: %v", err)
        http.Error(w, "error encoding JSON", http.StatusInternalServerError)
    }
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    session, err := ValidateSession(w, r, GetDB())
    if err != nil {
        log.Printf("Couldn't validate session: %v", err)
        http.Redirect(w, r, "/log-in", http.StatusSeeOther)
        return
    }

    // Fetch categories
    categories, err := GetAllCategories(GetDB())
    if err != nil {
        log.Printf("Error fetching categories: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    if r.Method == http.MethodPost {
        title := r.FormValue("title")
        content := r.FormValue("content")

        if content == "" {
            log.Printf("Content is empty")
            ErrorHandler(w, r, 400)
            return
        }

        category := r.Form["category"]
        if len(category) == 0 {
            log.Printf("No category selected")
            ErrorHandler(w, r, 400)
            return
        }

        if strings.Contains(content, "\n") {
            content = strings.ReplaceAll(content, "\n", "<br>")
        }

        // Create new post
        err = CreatePost(GetDB(), session.UserID, title, content)
        if err != nil {
            log.Printf("Error creating post: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Other logic here...

        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Prepare data for the template
    data := struct {
        Categories []Category
        Posts      []Posts
        Userid     int
        Post       Posts
        Comments   []Comments
        Commentusername string
    }{
        Categories: categories,
        Posts:      []Posts{},
        Userid:     0,
        Post:       Posts{},
        Comments:   []Comments{},
        Commentusername: "",
    }

    tmpl, err := template.ParseFiles("./index.html")
    if err != nil {
        log.Printf("Error parsing template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = tmpl.ExecuteTemplate(w, "index.html", data)
    if err != nil {
        log.Printf("Error executing template: %v", err)
        ErrorHandler(w, r, 500)
        return
    }
}


func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postIDStr := r.FormValue("post_id")
		content := r.FormValue("content")
		if content == "" {
			ErrorHandler(w, r, 400)
		}
		
		if content == "" {
			ErrorHandler(w, r, 400)
		}
		
		// Convert postID to int
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			log.Printf("Error converting post_id to int: %v", err)
			http.Error(w, "Invalid Post ID", http.StatusBadRequest)
			return
		}
		// Validate session
		session, err := ValidateSession(w, r, GetDB())
		if err != nil {
			log.Printf("Error validating session: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		username, ree := GetUsername(GetDB(), session.UserID)
		if ree != nil {
			log.Printf("Error fetching username: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
		if strings.Contains(content, "<br>") {
			content = strings.ReplaceAll(content, "<br>","<2br>")
		}
		// Create the comment
		err = CreateComment(GetDB(), session.UserID, postID, content, username)
		if err != nil {
			log.Printf("Error creating comment: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Redirect to the post page
		http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func PostCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Manually parse the URL to get the post ID
	if r.Method == "GET" {

	path := strings.Split(r.URL.Path, "/")
	id := path[2]
	postID, err := strconv.Atoi(string(id))
	if err != nil {
		log.Printf("Invalid post ID: %v", err)
		http.NotFound(w, r)
		return
	}
	// Fetch post details
	post, err := GetPostByID(GetDB(), postID)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Fetch comments for the post
	comments, err := GetCommentsByPostID(GetDB(), post.ID)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	post.Username, err = GetUsername(GetDB(), post.UserID)
	if err != nil {
		log.Printf("Error fetching username: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Display the post and its comments
	data := struct {
		Userid             int
		Post               Posts
		Comments           []Comments
		Commentusername    string
		UserID            int
		Posts             []Posts
		Categories []Category;
	}{
		Userid:             0,
		Post:               post,
		Comments:           comments,
		Commentusername:    "", // We'll set this value
		UserID:            0,
		Posts:             []Posts{},
		Categories:  []Category{},
		
	}
	
	for _, comment := range comments {
		comuid := comment.UserID
		Comusername, err := GetUsername(GetDB(), comuid)
		if err != nil {
			log.Printf("Error fetching username: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		data.Commentusername = Comusername
	}
	tmpl, err := template.ParseFiles("./index.html")
	if err != nil {
		ErrorHandler(w, r, 500)
		return
	}
	err = tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		ErrorHandler(w, r, 500)
		return
	}
}
}


func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
    /*
	session, err := ValidateSession(w, r, GetDB())
	if err != nil {
		return // The ValidateSession function handles redirection on error
	}
	*/
	posts, err := GetAllPosts(GetDB())
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		return
	}
	
	for i, post := range posts {
	
		// Fetch categories for each post
		
		categoriesID, err := GetCategoriesOfPosts(GetDB(), post.ID)
		if err != nil {
			log.Printf("Error fetching categories for post %d: %v", post.ID, err)
			continue
		}
		
		category, err := GetCategoryNamesByIDs(GetDB(), categoriesID)
		if err != nil {
			log.Printf("Error fetching category names for IDs %v: %v", categoriesID, err)
			continue
		}
		
		post.CategoriesNames = category
		if strings.Contains(post.Content, "\n") {
			post.Content = strings.ReplaceAll(post.Content, "\n", "\\n")
		
		}
	
		posts[i] = post
	}
	
}


func ServeJSON(w http.ResponseWriter, r *http.Request) {
	categories,err := GetAllCategories(GetDB())
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}

func ServeJSONPost(w http.ResponseWriter, r *http.Request) {
	posts, err := GetAllPosts(GetDB())
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}
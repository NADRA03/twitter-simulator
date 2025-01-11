package twitter

import(
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
    "time"
    "strings"
    "net/url"
)

type User struct {
	Id        int            `json:"id"`          // Unique identifier for the user
	Username  string         `json:"username"`    // Unique username for the user
	Email     string         `json:"email"`       // Unique email for the user
	FirstName string         `json:"first_name"`  // User's first name
	LastName  string         `json:"last_name"`   // User's last name
	Age       int            `json:"age"`         // User's age
	Gender    string         `json:"gender"`      // User's gender
	Password  string         `json:"password"`    // User's password (hashed ideally)
	ImageURL  sql.NullString  `json:"image_url"`   // Optional image URL for the user
    BigImageURL sql.NullString `json:"big_image_url"` // Optional big image URL for the user
	Role      string         `json:"role"`        // User's role (e.g., admin, user)
	CreatedAt time.Time      `json:"created_at"`  // Timestamp for when the user was created  
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/log-in" {
        http.NotFound(w, r)
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

            // Add error handling and pass messages to the front-end
            exists, err := LoginUser(db, username, password)
            if err != nil {
                log.Printf("Error logging user: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            if !exists {
                log.Println("Invalid username or password")
                // Pass error to the HTML template (e.g., show on the page)
                http.Redirect(w, r, "/log-in?error=Invalid+username+or+password&form=login", http.StatusSeeOther)
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


func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/log-in/create-account" {
        http.NotFound(w, r)
        return
    }

    if r.Method == http.MethodPost {
        err := r.ParseForm()
        if err != nil {
            http.Redirect(w, r, "/log-in?error=Failed+to+parse+form&form=signup", http.StatusSeeOther)
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
                http.Redirect(w, r, "/log-in?error="+url.QueryEscape(err.Error())+"&form=signup", http.StatusSeeOther)
                return
            }

            if err := ValidPassword(password); err != nil {
                http.Redirect(w, r, "/log-in?error="+url.QueryEscape(err.Error())+"&form=signup", http.StatusSeeOther)
                return
            }

            if email == "" {
                http.Redirect(w, r, "/log-in?error=Email+cannot+be+empty&form=signup", http.StatusSeeOther)
                return
            }

            // Check if username or email already exists
            exists, err := UserExists(db, username, email)
            if err != nil {
                log.Printf("Error checking user existence: %v", err)
                http.Redirect(w, r, "/log-in?error=Internal+Server+Error&form=signup", http.StatusSeeOther)
                return
            }

            if exists {
                http.Redirect(w, r, "/log-in?error=Username+or+email+already+exists&form=signup", http.StatusSeeOther)
                return
            }

            hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
            if err != nil {
                http.Redirect(w, r, "/log-in?error=Internal+Server+Error&form=signup", http.StatusSeeOther)
                return
            }

            err = UserRegister(db, username, email, firstname, lastname, age, gender, string(hashedPassword))
            if err != nil {
                http.Redirect(w, r, "/log-in?error=Internal+Server+Error&form=signup", http.StatusSeeOther)
                return
            }

            userID, _ := GetUserID(db, username)
            err = CreateSession(w, r, db, userID)
            if err != nil {
                log.Printf("Error creating session: %v", err)
                http.Redirect(w, r, "/log-in?error=Internal+Server+Error&form=signup", http.StatusSeeOther)
                return
            }

            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }
    }
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
    err = db.QueryRow("SELECT id, username, email, image_url, big_image_url FROM Users WHERE id = ?", session.UserID).
        Scan(&user.Id, &user.Username, &user.Email, &user.ImageURL, &user.BigImageURL)
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

    // Check if ImageURL and BigImageURL are valid, and set them to empty strings if NULL
    if !user.ImageURL.Valid {
        user.ImageURL.String = ""
    }
    if !user.BigImageURL.Valid {
        user.BigImageURL.String = ""
    }

    log.Printf("UserDetailsHandler: User details retrieved: %+v", user)

    // Send the user details as JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user); err != nil {
        log.Printf("UserDetailsHandler: Error encoding user data to JSON: %v", err)
        http.Error(w, "error encoding JSON", http.StatusInternalServerError)
    }
}


func AUserDetailsHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("UserDetailsHandler: Profile is loading")

    // Extract the userID from the URL path
    path := r.URL.Path
    // Split the URL path to get the user ID (assuming the URL structure is /a_profile/{id})
    parts := strings.Split(path, "/")

    // Ensure there are at least 3 parts in the path: ["", "a_profile", "{id}"]
    if len(parts) < 3 {
        log.Println("UserDetailsHandler: Invalid URL path")
        http.Error(w, "invalid URL", http.StatusBadRequest)
        return
    }

    // The ID should be the last part of the path
    userID := parts[len(parts)-1]

    // Check if userID is valid
    if userID == "" {
        log.Println("UserDetailsHandler: No user ID provided in URL path")
        http.Error(w, "user ID is required", http.StatusBadRequest)
        return
    }

    log.Printf("UserDetailsHandler: Retrieving details for UserID %s", userID)

    // Query the user details based on the userID from the URL path
    var user User
    err := db.QueryRow("SELECT id, username, email, image_url, big_image_url FROM Users WHERE id = ?", userID).
        Scan(&user.Id, &user.Username, &user.Email, &user.ImageURL, &user.BigImageURL)
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

    // Check if ImageURL and BigImageURL are valid, and set them to empty strings if NULL
    if !user.ImageURL.Valid {
        user.ImageURL.String = ""
    }
    if !user.BigImageURL.Valid {
        user.BigImageURL.String = ""
    }

    log.Printf("UserDetailsHandler: User details retrieved: %+v", user)

    // Send the user details as JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(user); err != nil {
        log.Printf("UserDetailsHandler: Error encoding user data to JSON: %v", err)
        http.Error(w, "error encoding JSON", http.StatusInternalServerError)
    }
}

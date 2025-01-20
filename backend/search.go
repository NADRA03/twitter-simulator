package twitter

import(
	"net/http"
	"log"
	"encoding/json"
)


func SearchUsers(w http.ResponseWriter, r *http.Request) {
    // log.Println("Request received to search for users")

    // Get the search term from the query parameters
    searchTerm := r.URL.Query().Get("term")
    if searchTerm == "" {
        http.Error(w, "Search term is required", http.StatusBadRequest)
        return
    }

    // Prepare the SQL query to search for users
    query := `
        SELECT id, username, image_url
        FROM users
        WHERE username LIKE ?
        LIMIT 10
    `

    // Execute the query
    rows, err := db.Query(query, "%"+searchTerm+"%")
    if err != nil {
        log.Printf("Error executing query: %v", err)
        http.Error(w, "Error fetching users", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Initialize a slice to hold the matched users
    var matchedUsers []User

    // Scan the results
    for rows.Next() {
        var user User
        err := rows.Scan(&user.Id, &user.Username, &user.ImageURL) // Update scan to include ID
        if err != nil {
            log.Printf("Error scanning user data: %v", err)
            http.Error(w, "Error scanning user data", http.StatusInternalServerError)
            return
        }
        matchedUsers = append(matchedUsers, user)
    }

    // Check for errors that occurred during iteration
    if err := rows.Err(); err != nil {
        log.Printf("Error reading user data: %v", err)
        http.Error(w, "Error reading user data", http.StatusInternalServerError)
        return
    }

    // Set the response header to JSON and send the matched users
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(matchedUsers); err != nil {
        log.Printf("Error encoding response as JSON: %v", err)
        http.Error(w, "Error encoding response as JSON", http.StatusInternalServerError)
    }
}

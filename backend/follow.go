package twitter

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"encoding/json"
)

func AddFollowHandler(w http.ResponseWriter, r *http.Request) {
	// Validate session to get userID
	session, err := ValidateSession(w, r, db)
	if err != nil {
		log.Println("addFollowHandler - ValidateSession: Session validation failed", err)
		http.Error(w, "Invalid session", http.StatusForbidden)
		return
	}
	userID := session.UserID
	fmt.Printf("Session validated, userID: %d\n", userID)

	// Get followed_id from the URL query parameters
	followedIDStr := r.URL.Query().Get("followed_id")
	followedID, err := strconv.Atoi(followedIDStr)
	if err != nil || followedID <= 0 {
		http.Error(w, "Invalid followed ID", http.StatusBadRequest)
		return
	}

	// Insert the follow relationship into the database
	_, err = db.Exec("INSERT OR IGNORE INTO follow (follower_id, followed_id) VALUES (?, ?)", userID, followedID)
	if err != nil {
		http.Error(w, "Could not follow user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Follow added successfully"))
}

func FollowingCountHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL query parameter
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("Fetching following count for userID: %d\n", userID)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM follow WHERE follower_id = ?", userID).Scan(&count)
	if err != nil {
		http.Error(w, "Could not get following count", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Following count: %d", count)))
}

// Handler to get the followers count for a user by userID from the request parameter
func FollowersCountHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL query parameter
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("Fetching followers count for userID: %d\n", userID)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM follow WHERE followed_id = ?", userID).Scan(&count)
	if err != nil {
		http.Error(w, "Could not get followers count", http.StatusInternalServerError)
		return
	}

	// Send the followers count as JSON
	response := map[string]int{"followersCount": count}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
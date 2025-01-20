package twitter

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"encoding/json"
)

func AddFollowHandler(w http.ResponseWriter, r *http.Request) {
	session, err := ValidateSession(w, r, db)
	if err != nil {
		log.Println("addFollowHandler - ValidateSession: Session validation failed", err)
		http.Error(w, "Invalid session", http.StatusForbidden)
		return
	}
	userID := session.UserID
	fmt.Printf("Session validated, userID: %d\n", userID)

	followedIDStr := r.URL.Query().Get("followed_id")
	followedID, err := strconv.Atoi(followedIDStr)
	if err != nil || followedID <= 0 {
		http.Error(w, "Invalid followed ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT OR IGNORE INTO follow (follower_id, followed_id) VALUES (?, ?)", userID, followedID)
	if err != nil {
		http.Error(w, "Could not follow user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Follow added successfully"))
}

func FollowingCountHandler(w http.ResponseWriter, r *http.Request) {
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


func FollowersCountHandler(w http.ResponseWriter, r *http.Request) {

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

	response := map[string]int{"followersCount": count}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func TopFollowedUserIDsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := ValidateSession(w, r, db)
	if err != nil {
		log.Println("TopFollowedUserIDsHandler - ValidateSession: Session validation failed", err)
		http.Error(w, "Invalid session", http.StatusForbidden)
		return
	}
	userID := session.UserID
	fmt.Printf("Session validated, userID: %d\n", userID)

	query := `
        SELECT u.id
        FROM users u
        JOIN follow f ON u.id = f.followed_id
        WHERE u.id != ? AND NOT EXISTS (
            SELECT 1 FROM follow existing_follow WHERE existing_follow.follower_id = ? AND existing_follow.followed_id = u.id
        )
        GROUP BY u.id
        ORDER BY COUNT(f.follower_id) DESC
        LIMIT 10`

	rows, err := db.Query(query, userID, userID)
	if err != nil {
		log.Printf("TopFollowedUserIDsHandler: Error querying top followed user IDs: %v", err)
		http.Error(w, "Could not get top followed user IDs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var topUserIDs []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("TopFollowedUserIDsHandler: Error scanning row: %v", err)
			http.Error(w, "Failed to retrieve user IDs", http.StatusInternalServerError)
			return
		}
		topUserIDs = append(topUserIDs, id)
	}

	if err = rows.Err(); err != nil {
		log.Printf("TopFollowedUserIDsHandler: Error during row iteration: %v", err)
		http.Error(w, "Failed during row iteration", http.StatusInternalServerError)
		return
	}

	if len(topUserIDs) == 0 {
		// fmt.Println("No top followed users found.")
	} else {
		fmt.Printf("Top followed user IDs: %v\n", topUserIDs)
	}


	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(topUserIDs); err != nil {
		log.Printf("TopFollowedUserIDsHandler: Error encoding JSON response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
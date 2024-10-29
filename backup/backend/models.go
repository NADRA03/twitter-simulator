package backend

import(
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func UserExists(db *sql.DB, username, email string) (bool, error) {
	var count int 
	query := `SELECT COUNT(*) FROM users WHERE username = ? OR email = ?`
	err := db.QueryRow(query,username,email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count >0 , nil
}

func ValidUsername(db *sql.DB, username string) error {
	if len(username) == 0 {
		return fmt.Errorf("user name cannot be empty")
	}

	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 charecters")
	}

	if len(username) > 25 {
		return fmt.Errorf("username must be less than 25 charecters")
	}

	matched, err := regexp.MatchString("^[a-zA-Z0-9_.]+$",username)
	if err != nil {
		return fmt.Errorf("error validating username : %v", err)
	}
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and dots")
	}

	if strings.Contains(username," ") {
		return fmt.Errorf("username cannot contain spaces")
	}

	// Check if the username is already taken
	err = db.QueryRow("SELECT 1 FROM users WHERE nickname = ?", username).Scan(&[]byte{})
	if err == nil {
		fmt.Errorf("username already exists")
	} else if err != sql.ErrNoRows {
		fmt.Errorf("error querying database: %v", err)
	}
	return nil
}


func ValidPassword(password string) error {
	if len(password) == 0 {
		return fmt.Errorf("password cannot be empty")
	}

	if len(password) < 6 {
		return fmt.Errorf("password cannot be less than 6 characters")
	}

	if len(password) > 15 {
		return fmt.Errorf("password cannot be more than 15 characters")
	}

	matched, err := regexp.MatchString("[A-Z]",password)
	if err != nil {
	 	return fmt.Errorf("error validating password :%v", err)
	}
	if !matched {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	matched, err = regexp.MatchString("[a-z]", password)
	if err != nil {
		return fmt.Errorf("error validating password: %v", err)
	}
	if !matched {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	
	matched, err = regexp.MatchString("[0-9]", password)
	if err != nil {
		return fmt.Errorf("error validating password: %v", err)
	}
	if !matched {
		return fmt.Errorf("password must contain at least one digit")
	}
	
	if matched, _ := regexp.MatchString(`\s`, password); matched {
		return fmt.Errorf("password cannot contain spaces")
	}
	// Check if the password contains any symbols
	matched, err = regexp.MatchString(`[!@#\$%\^&\*(),.?":{}|<>]`, password)
	if err != nil {
		return fmt.Errorf("error validating password: %v", err)
	}
	if matched {
		return fmt.Errorf("password cannot contain special symbols")
	}
	return nil
}

func IsUserLoggedIn(db *sql.DB, username string)(bool, error) {
	var userID int 
	err := db.QueryRow("SELECT id FROM Sessions WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error querying user ID: %v", err)
	}
	return true, nil
}

func LoginUser(db *sql.DB, username, password string,) (bool, error) {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ? OR email = ?", username, username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Username not found
			fmt.Println("Username not found")
			return false, nil
		}
		return false, fmt.Errorf("error querying database: %v", err)
	}

	// Compare the provided password with the hashed password from the database
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Invalid password
		fmt.Println("Invalid password")
		return false, nil
	}

	return true,nil
}

func GetUsername(db *sql.DB, userID int) (string, error) {
	var username string
	err := db.QueryRow("SELECT nickname FROM Users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user ID not found")
		}
		return "", fmt.Errorf("error querying username: %v", err)
	}
	return username, nil
}

func GetUserID(db *sql.DB, username string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE nickname = ? OR email = ?", username, username).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error querying user ID: %v", err)
	}
	return userID, nil
}
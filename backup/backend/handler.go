package backend

import(
	"net/http"
	"database/sql"
	//"encoding/json"
	"text/template"
	"fmt"
	"log"
	
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    
    temp, err := template.ParseFiles("frontend/index.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    err = temp.ExecuteTemplate(w, "index.html", nil)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

func LoginHandeler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w,r)
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
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Login successful"))
				return
			}

		}


	} else {
		tmpl, err := template.ParseFiles("frontend/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		
		err = tmpl.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	}

}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/signup" {
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

            http.Redirect(w, r, "/home", http.StatusSeeOther)
            return
        }
    }

    tmpl, err := template.ParseFiles("frontend/index.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = tmpl.ExecuteTemplate(w, "index.html", nil)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}



func HomePageHandler(w http.ResponseWriter, r *http.Request) {
    
	if r.Method != http.MethodGet { // Check for GET method
        http.NotFound(w, r)
        return
    }

    /*
    session, err := ValidateSession(w, r, GetDB())
    if err != nil {
        return 
    }
*/
   
    tmpl, err := template.ParseFiles("frontend/index.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Execute the template
    err = tmpl.ExecuteTemplate(w, "index.html", nil)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}


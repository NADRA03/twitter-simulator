package twitter

import (
	"io"
	"net/http"
	"text/template"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	if status == http.StatusInternalServerError {
		t, err := template.ParseFiles("frontend/html/internal.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "500\nThere is an internal server error")
			return 
		}
		w.WriteHeader(status)
		t.Execute(w, "frontend/html/internal.html")
	}
	if status == http.StatusNotFound {
		t, err := template.ParseFiles("frontend/html/notfound.html")
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(status)
		t.Execute(w, "frontend/html/notfound.html")
	} else if status == http.StatusBadRequest {
		t, err := template.ParseFiles("frontend/html/badrequest.html")
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError)
			// w.WriteHeader(http.StatusInternalServerError)
			// t = template.Must(template.ParseFiles("templates/error.html"))
			return
		}
		w.WriteHeader(status)
		t.Execute(w, "frontend/html/badrequest.html")
	}
}
package restAPI

import "net/http"

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotAcceptable)
	w.Write([]byte("{Message : 'Wrong URL/Method'}"))
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		ErrorHandler(w, r)
		return
	}
	w.Write([]byte("This is the HOME PAGE"))
}

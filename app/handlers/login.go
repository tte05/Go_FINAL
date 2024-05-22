package handlers

import (
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"goproject/app/db"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("advanced"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/login.html")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "advanced")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = nil
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := db.FindUserByUsername(username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !user.Confirmed {
		http.Error(w, "You need to confirm your account before logging in", http.StatusUnauthorized)
		return
	}

	session, err := store.Get(r, "advanced")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = user.ID.Hex()
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

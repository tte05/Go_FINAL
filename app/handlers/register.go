package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"goproject/app/db"
	"goproject/app/models"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/register.html")
}
func RegisterSubmitHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	confirmationToken := generateToken()

	newUser := models.User{
		Username:          username,
		Password:          string(hashedPassword),
		Email:             email,
		Confirmed:         false,
		ConfirmationToken: confirmationToken,
		Role:              "user",
	}

	err = db.InsertUser(newUser)
	if err != nil {
		if err.Error() == "username already exists" {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}
		if err.Error() == "email already exists" {
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	err = SendConfirmationEmail(newUser)
	if err != nil {
		http.Error(w, "Failed to send confirmation email", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

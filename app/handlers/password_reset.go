package handlers

import (
	"golang.org/x/crypto/bcrypt"
	"goproject/app/db"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func ShowPasswordResetRequestForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/password_reset_request.html")
}

func HandlePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	user, err := db.FindUserByEmail(email)
	if err != nil {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}

	token := uuid.New().String()
	user.ResetToken = token
	user.ResetTokenExpiry = time.Now().Add(1 * time.Hour)
	err = db.UpdateUser(user)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	err = SendPasswordResetEmail(user.Email, token)
	if err != nil {
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ShowPasswordResetForm(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	user, err := db.FindUserByResetToken(token)
	if err != nil || time.Now().After(user.ResetTokenExpiry) {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	t, _ := template.ParseFiles("templates/password_reset.html")
	t.Execute(w, struct{ Token string }{Token: token})
}

func HandlePasswordReset(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	newPassword := r.FormValue("new-password")

	user, err := db.FindUserByResetToken(token)
	if err != nil || time.Now().After(user.ResetTokenExpiry) {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}
	err = db.UpdateUser(user)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

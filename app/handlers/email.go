package handlers

import (
	"github.com/go-gomail/gomail"
	"github.com/gorilla/mux"
	"goproject/app/db"
	"goproject/app/models"
	"net/http"
)

func SendConfirmationEmail(user models.User) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "javafrom@yandex.ru")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Confirmation of registration")
	m.SetBody("text/html", "To confirm registration, click on the link: <a href='https://go-final-1n5l.onrender.com/confirm/"+user.ConfirmationToken+"'>Confirm</a>")

	d := gomail.NewDialer("smtp.yandex.ru", 587, "javafrom@yandex.ru", "cegbkuthamcmvsvm")

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]

	user, err := db.FindUserByConfirmationToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	user.Confirmed = true
	user.ConfirmationToken = ""
	err = db.UpdateUser(user)
	if err != nil {
		http.Error(w, "Failed to confirm email", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, "templates/email_confirmed.html")
}
func SendPasswordResetEmail(to string, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "javafrom@yandex.ru")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/html", "Please use the following link to reset your password: <a href='https://go-final-1n5l.onrender.com/password-reset?token="+token+"'>Reset</a>")

	d := gomail.NewDialer("smtp.yandex.ru", 587, "javafrom@yandex.ru", "cegbkuthamcmvsvm")

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

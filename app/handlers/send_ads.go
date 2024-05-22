package handlers

import (
	"gopkg.in/gomail.v2"
	"goproject/app/db"
	"net/http"
)

func SendAdsSubmitHandler(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	users, err := db.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		err := sendAdEmail(user.Email, message)
		if err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func sendAdEmail(email, message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "javafrom@yandex.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Advertisement")
	m.SetBody("text/html", message)

	d := gomail.NewDialer("smtp.yandex.ru", 587, "javafrom@yandex.ru", "cegbkuthamcmvsvm")

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

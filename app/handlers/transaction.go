package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"goproject/app/db"
	"goproject/app/models"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BuyPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/payment.html")
}
func BuySuccess(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/payment_s.html")
}
func PaymentSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	game := r.FormValue("gameName")
	email := r.FormValue("email")

	paymentSuccessful := true

	if paymentSuccessful {
		transaction := models.Transaction{
			Project: game,
			Customer: models.Customer{
				Name:  name,
				Email: email,
			},
			Status: "Paid",
			Date:   time.Now(),
		}

		transaction.ID = primitive.NewObjectID()
		err := db.CreateTransaction(&transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = GenerateAndSendReceipt(transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/success", http.StatusSeeOther)
	} else {
		http.Error(w, "Payment failed", http.StatusInternalServerError)
	}
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction.ID = primitive.NewObjectID()
	transaction.Date = time.Now()

	err = db.CreateTransaction(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	paymentSuccessful := true

	if paymentSuccessful {
		update := bson.M{"$set": bson.M{"status": "Paid"}}
		err := db.UpdateTransactionByID(id, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		transaction := models.Transaction{
			ID:       id,
			Project:  "Project Name",
			Customer: models.Customer{Name: "Alibi", Email: "220940@astanait.edu.kz"},
			Status:   "Paid",
		}
		err = GenerateAndSendReceipt(transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func GenerateAndSendReceipt(transaction models.Transaction) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Fiscal Receipt")

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Project: %s", transaction.Project))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Transaction ID: %s", transaction.ID.Hex()))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", transaction.Date.Format("2006-01-02 15:04:05")))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", transaction.Customer.Name))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Price: 10$"))
	pdf.Ln(10)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return err
	}

	err = sendEmail(transaction.Customer.Email, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func sendEmail(to string, attachment []byte) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "javafrom@yandex.ru")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Fiscal Receipt")
	m.SetBody("text/html", "Please find attached your fiscal receipt.")

	m.Attach("receipt.pdf", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(attachment)
		return err
	}))

	d := gomail.NewDialer("smtp.yandex.ru", 587, "javafrom@yandex.ru", "cegbkuthamcmvsvm")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	fmt.Println("gtrer")
	return nil
}

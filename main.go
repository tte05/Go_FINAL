package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"goproject/app/db"
	"goproject/app/middleware"
	"goproject/app/router"
	"net/http"
	"os"
)

const logFilePath = "app.log"

var logger *logrus.Logger

func main() {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Unable to create log file: " + err.Error())
	}
	defer logFile.Close()

	logger = logrus.New()
	logger.Out = logFile

	client, err := db.ConnectToMongoDB()
	if err != nil {
		logger.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	router := router.SetupRouter()
	limiter := middleware.NewLimiter(2, 5)
	router.Use(limiter)
	
	router.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.html")
	}).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/games", http.StatusSeeOther)
	})

	http.Handle("/", router)

	logger.Info("Server started at localhost:8080")
	logger.Fatal(http.ListenAndServe(":8080", router))
}

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
	logger = logrus.New()

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("Unable to create log file: %v", err)
	}
	defer logFile.Close()

	logger.Out = logFile

	client, err := db.ConnectToMongoDB()
	if err != nil {
		logger.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	r := router.SetupRouter()
	limiter := middleware.NewLimiter(2, 5)
	r.Use(limiter)

	r.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.html")
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/games", http.StatusSeeOther)
	})

	http.Handle("/", r)

	port := os.Getenv("PORT")

	logger.Infof("Server started at localhost:%s", port)

	logger.Fatal(http.ListenAndServe(":"+port, nil))
}

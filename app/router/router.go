package router

import (
	"github.com/gorilla/mux"
	"goproject/app/handlers"
	"goproject/app/middleware"
	"net/http"
)

func SetupRouter() *mux.Router {
	store := middleware.InitSessionStore()

	router := mux.NewRouter()

	router.HandleFunc("/register", handlers.RegisterHandler).Methods("GET")
	router.HandleFunc("/register", handlers.RegisterSubmitHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("GET")
	router.HandleFunc("/login", handlers.LoginSubmitHandler).Methods("POST")
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")

	router.HandleFunc("/confirm/{token}", handlers.ConfirmEmailHandler).Methods("GET")

	router.HandleFunc("/password-reset-request", handlers.ShowPasswordResetRequestForm).Methods("GET")
	router.HandleFunc("/password-reset-request", handlers.HandlePasswordResetRequest).Methods("POST")
	router.HandleFunc("/password-reset", handlers.ShowPasswordResetForm).Methods("GET")
	router.HandleFunc("/password-reset", handlers.HandlePasswordReset).Methods("POST")

	auth := router.PathPrefix("/").Subrouter()
	auth.Use(middleware.AuthRequired(store))

	auth.HandleFunc("/games", handlers.GetGamesHandler).Methods("GET")
	auth.HandleFunc("/games", handlers.CreateGameHandler).Methods("POST")
	auth.HandleFunc("/games/{id}", handlers.GetGameByIDHandler).Methods("GET")
	auth.HandleFunc("/games/{id}", handlers.UpdateGameByIDHandler).Methods("PUT")
	auth.HandleFunc("/games/{id}", handlers.DeleteGameByIDHandler).Methods("DELETE")

	auth.HandleFunc("/buy", handlers.BuyPageHandler).Methods("GET")
	auth.HandleFunc("/success", handlers.BuySuccess).Methods("GET")
	auth.HandleFunc("/transactions/pay", handlers.PaymentSubmitHandler).Methods("POST")
	auth.HandleFunc("/transactions", handlers.CreateTransactionHandler).Methods("POST")
	auth.HandleFunc("/transactions/{id}/pay", handlers.PaymentHandler).Methods("POST")

	admin := auth.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminRequired(store))
	admin.HandleFunc("", handlers.AdminPageHandler).Methods("GET")
	admin.HandleFunc("", handlers.SendAdsSubmitHandler).Methods("POST")
	admin.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")
	admin.HandleFunc("/users/{id}", handlers.DeleteUserByIDHandler).Methods("DELETE")
	admin.HandleFunc("/users/{id}/role", handlers.ChangeUserRoleHandler).Methods("PATCH")

	auth.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates")))

	return router
}

package middleware

import (
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goproject/app/db"
	"net/http"
)

var store *sessions.CookieStore

func InitSessionStore() *sessions.CookieStore {
	store = sessions.NewCookieStore([]byte("advanced"))
	return store
}

func AuthRequired(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "advanced")
			if session.Values["userID"] == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AdminRequired(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "advanced")
			userID := session.Values["userID"]
			if userID == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			objID, _ := primitive.ObjectIDFromHex(userID.(string))
			user, err := db.FindUserByID(objID)
			if err != nil || user.Role != "admin" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

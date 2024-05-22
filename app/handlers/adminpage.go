package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goproject/app/db"
	"goproject/app/models"
	"net/http"
)

func AdminPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/adminpage.html")
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	cursor, err := db.Collection2.Find(context.Background(), bson.M{})
	if err != nil {
		handleError(w, err)
		return
	}
	defer cursor.Close(context.Background())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			handleError(w, err)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"users": users}); err != nil {
		handleError(w, err)
		return
	}
}

func DeleteUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handleError(w, err)
		return
	}
	_, err = db.Collection2.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func ChangeUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handleError(w, err)
		return
	}

	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		handleError(w, err)
		return
	}

	newRole, ok := requestBody["role"]
	if !ok {
		http.Error(w, "role field is required", http.StatusBadRequest)
		return
	}

	_, err = db.Collection2.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"role": newRole}})
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User role updated successfully"})
}

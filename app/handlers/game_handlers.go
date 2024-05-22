package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goproject/app/db"
	"goproject/app/models"
	"math"
	"net/http"
	"strconv"
)

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = db.CreateGame(&game)
	if err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Game created successfully"})
}

func GetGamesHandler(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 5
	}

	offset := (page - 1) * pageSize

	sortOptions := bson.M{}
	switch sortBy {
	case "titleASC":
		sortOptions["title"] = 1
	case "genreASC":
		sortOptions["genre"] = 1
	case "ratingASC":
		sortOptions["rating"] = 1
	case "titleDESC":
		sortOptions["title"] = -1
	case "genreDESC":
		sortOptions["genre"] = -1
	case "ratingDESC":
		sortOptions["rating"] = -1
	}

	options := options.Find().SetSort(sortOptions).SetSkip(int64(offset)).SetLimit(int64(pageSize))
	filter := bson.M{}
	minRating, err := strconv.ParseFloat(r.URL.Query().Get("minRating"), 64)
	if err == nil {
		filter["rating"] = bson.M{"$gte": minRating}
	}
	cursor, err := db.Collection.Find(context.Background(), filter, options)

	if err != nil {
		handleError(w, err)
		return
	}
	defer cursor.Close(context.Background())

	var games []models.Game
	for cursor.Next(context.Background()) {
		var game models.Game
		if err := cursor.Decode(&game); err != nil {
			handleError(w, err)
			return
		}
		games = append(games, game)
	}

	totalGames, err := db.Collection.CountDocuments(context.Background(), filter)
	if err != nil {
		handleError(w, err)
		return
	}
	totalPages := int(math.Ceil(float64(totalGames) / float64(pageSize)))

	responseData := map[string]interface{}{
		"games":      games,
		"totalPages": totalPages,
		"sortBy":     sortBy,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		handleError(w, err)
		return
	}
}

func GetGameByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handleError(w, err)
		return
	}
	var game models.Game
	err = db.Collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&game)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(game)
}

func UpdateGameByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]
	if gameID == "" {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		http.Error(w, "the provided hex string is not a valid ObjectID", http.StatusBadRequest)
		return
	}

	var game models.Game
	err = json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	game.ID = objectID
	err = db.UpdateGameByID(&game)
	if err != nil {
		http.Error(w, "Failed to update game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Game updated successfully"})
}

func DeleteGameByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]
	if gameID == "" {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		http.Error(w, "the provided hex string is not a valid ObjectID", http.StatusBadRequest)
		return
	}

	err = db.DeleteGameByID(objectID)
	if err != nil {
		http.Error(w, "Failed to delete game", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Game deleted successfully"})
}

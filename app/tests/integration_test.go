package tests

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goproject/app/db"
	"goproject/app/handlers"
	"goproject/app/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testClient *mongo.Client
var testDB *mongo.Database

func TestMain(m *testing.M) {
	var err error
	testClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	err = testClient.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	testDB = testClient.Database("gamesdb_test")
	db.Collection = testDB.Collection("games")
	code := m.Run()

	testDB.Drop(context.Background())
	testClient.Disconnect(context.Background())

	os.Exit(code)
}

func TestGetGamesHandlerIntegration(t *testing.T) {
	game := models.Game{
		Title:       "Test Game",
		Genre:       "Adventure",
		Rating:      5,
		Developer:   "Test Developer",
		Description: "A test game description",
	}
	_, err := db.Collection.InsertOne(context.Background(), game)
	if err != nil {
		t.Fatalf("Failed to insert test game: %v", err)
	}
	req, _ := http.NewRequest("GET", "/games", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetGamesHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	if len(response["games"].([]interface{})) != 1 {
		t.Errorf("handler returned unexpected number of games: got %v want %v", len(response["games"].([]interface{})), 1)
	}

	// Clean up the test game
	_, err = db.Collection.DeleteOne(context.Background(), bson.M{"title": "Test Game"})
	if err != nil {
		t.Fatalf("Failed to delete test game: %v", err)
	}
}

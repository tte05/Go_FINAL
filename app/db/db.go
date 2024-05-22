package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goproject/app/models"
	"time"
)

const (
	connectionString = "mongodb+srv://Kiraro:EvMgE22srSkbb5Lc@cluster0.gkbbkh8.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	dbName           = "gamesdb"
)

var Collection *mongo.Collection
var Collection2 *mongo.Collection

func ConnectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	Collection = client.Database(dbName).Collection("games")
	Collection2 = client.Database(dbName).Collection("users")
	return client, nil
}

func InsertUser(user models.User) error {
	exists, err := UsernameExists(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already exists")
	}

	exists, err = EmailExists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	_, err = Collection2.InsertOne(context.TODO(), user)
	return err
}

func UsernameExists(username string) (bool, error) {
	filter := bson.M{"username": username}
	count, err := Collection2.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func EmailExists(email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := Collection2.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func FindUserByUsername(username string) (models.User, error) {
	var user models.User
	filter := bson.M{"username": username}
	err := Collection2.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func FindUserByEmail(email string) (models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := Collection2.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func FindUserByResetToken(token string) (models.User, error) {
	var user models.User
	filter := bson.M{"reset_token": token}
	err := Collection2.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func FindUserByConfirmationToken(token string) (models.User, error) {
	var user models.User
	filter := bson.M{"confirmation_token": token}
	err := Collection2.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func UpdateUser(user models.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"username":           user.Username,
			"password":           user.Password,
			"email":              user.Email,
			"confirmed":          user.Confirmed,
			"confirmation_token": user.ConfirmationToken,
			"reset_token":        user.ResetToken,
			"reset_token_expiry": user.ResetTokenExpiry,
		},
	}
	_, err := Collection2.UpdateOne(context.Background(), filter, update)
	return err
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	cursor, err := Collection2.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func FindUserByID(id primitive.ObjectID) (models.User, error) {
	var user models.User
	err := Collection2.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	return user, err
}

func CreateGame(game *models.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Collection.InsertOne(ctx, game)
	return err
}
func DeleteGameByID(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func UpdateGameByID(game *models.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":       game.Title,
			"genre":       game.Genre,
			"rating":      game.Rating,
			"developer":   game.Developer,
			"description": game.Description,
		},
	}

	_, err := Collection.UpdateOne(ctx, bson.M{"_id": game.ID}, update)
	return err
}

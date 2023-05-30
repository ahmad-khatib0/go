package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	model "github.com/ahmad-khatib0/go/basics/mongo/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "watchlist"

var collection *mongo.Collection

// init is a name comming from go like the main function
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var mongodbURL string = "mongodb+srv://" + os.Getenv(
		"MONGODB_USERNAME",
	) + ":" + os.Getenv(
		"MONGODB_PASSWORD",
	) + "@" + os.Getenv(
		"MONGODB_URL",
	) + "/" + os.Getenv(
		"MONGODB_DBNAME",
	) + "?retryWrites=true&w=majority"

	fmt.Println("mongo url is: ", mongodbURL)
	clientOption := options.Client().ApplyURI(mongodbURL)

	// this context is used to customize the behavior when connect to other machine, (there is background also, not just TODO)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connection is established successfully")
	collection = client.Database(os.Getenv("MONGODB_DBNAME")).Collection(collectionName)

	fmt.Println("collection reference is ready")
}

func insertOne(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted one record in db , its id is: ", inserted.InsertedID)
}

func updateOne(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)

	filter := bson.M{"_id": id} // mongo is using something like json but it called bson, its the same but more advanced
	update := bson.M{"$set": bson.M{"watched": true}}

	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", res.ModifiedCount)
}

func deleteOne(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)

	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", deleteCount)
}

func deleteAll() int64 {
	res, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("number of deleted movies: ", res.DeletedCount)
	return res.DeletedCount
}

func getAll() []primitive.M {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M
	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}

		movies = append(movies, movie)
	}
	defer cursor.Close(context.Background())
	return movies
}

// routes

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allMovies := getAll()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)

	insertOne(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOne(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOne(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deleteAll()
	json.NewEncoder(w).Encode(count)
}

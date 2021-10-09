package restAPI

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption  string             `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageUrl string             `json:"imageUrl,omitempty" bson:"imageUrl,omitempty"`
	Time     string             `json:"time,omitempty" bson:"time,omitempty"`
	UserId   string             `json:"userId,omitempty" bson:"userId"`
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	if r.Method != "POST" || r.URL.Path != "/posts" {
		ErrorHandler(w, r)
		return
	}
	w.Header().Set("content-type", "application/json")
	var post Post
	json.NewDecoder(r.Body).Decode(&post)
	collection := client.Database(DATABASE).Collection("posts")
	result, _ := collection.InsertOne(context.Background(), post)
	collection.FindOne(context.Background(), bson.D{{"_id", result.InsertedID}}).Decode(&post)
	json.NewEncoder(w).Encode(post)
	defer lock.Unlock()
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	pathElements := strings.Split(r.URL.Path, "/")
	if r.Method != "GET" || len(pathElements) != 3 {
		ErrorHandler(w, r)
		return
	}
	w.Header().Set("content-type", "application/json")
	id := pathElements[2]
	collection := client.Database(DATABASE).Collection("posts")
	postId, _ := primitive.ObjectIDFromHex(id)
	var post Post
	collection.FindOne(context.Background(), bson.D{{"_id", postId}}).Decode(&post)
	json.NewEncoder(w).Encode(post)
	defer lock.Unlock()
}

func GetPostsByUserHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	var results []Post
	var perPage int64 = 5
	page := 1

	pathElements := strings.Split(r.URL.Path, "/")
	w.Header().Set("content-type", "application/json")
	id := pathElements[3]
	if len(pathElements) == 5 {
		pageNo, _ := strconv.Atoi(pathElements[4])
		page = pageNo

	}

	findOptions := options.Find()
	findOptions.SetLimit(perPage)
	findOptions.SetSkip(int64((page - 1) * 5))
	collection := client.Database(DATABASE).Collection("posts")
	cursor, err := collection.Find(context.TODO(), bson.D{{"userId", id}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.Background()) {

		var post Post
		err := cursor.Decode(&post)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, post)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	cursor.Close(context.Background())
	json.NewEncoder(w).Encode(results)
	defer lock.Unlock()
}

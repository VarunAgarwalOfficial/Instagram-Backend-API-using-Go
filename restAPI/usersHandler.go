package restAPI

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omit"`
}

func CreateUsersHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	if r.Method != "POST" || r.URL.Path != "/users" {
		ErrorHandler(w, r)
		return
	}
	w.Header().Set("content-type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	encryptedPassword := md5.Sum([]byte(user.Password))
	user.Password = base64.StdEncoding.EncodeToString(encryptedPassword[:])
	collection := client.Database(DATABASE).Collection("users")
	result, _ := collection.InsertOne(context.Background(), user)
	collection.FindOne(context.Background(), bson.D{{"_id", result.InsertedID}}).Decode(&user)
	w.Write([]byte(`{
	_id :` + user.ID.Hex() + ` , 
	name: ` + user.Name + ` ,
	email: ` + user.Email + `}`))
	defer lock.Unlock()

}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	pathElements := strings.Split(r.URL.Path, "/")
	if r.Method != "GET" || len(pathElements) != 3 {
		ErrorHandler(w, r)
		return
	}
	w.Header().Set("content-type", "application/json")
	id := pathElements[2]
	collection := client.Database(DATABASE).Collection("users")
	userId, _ := primitive.ObjectIDFromHex(id)
	var user User
	collection.FindOne(context.Background(), bson.D{{"_id", userId}}).Decode(&user)
	json.NewEncoder(w).Encode(user)
	defer lock.Unlock()
}

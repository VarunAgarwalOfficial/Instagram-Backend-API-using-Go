package restAPI

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	MONGOURL = "mongodb://localhost:27017"
	DATABASE = "appointy"
	lock     sync.Mutex
)

func InitMongoConnection() {
	client, _ = mongo.NewClient(options.Client().ApplyURI(MONGOURL))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoServer :", MONGOURL, "has been successfully established")
}

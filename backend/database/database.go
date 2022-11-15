package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Clinet variable holds db instance and is accessable to other files
var Client *mongo.Client

func DBinstance() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env variables ", err)
	}

	MongoDBURL := os.Getenv("MONGODB_URL")
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(MongoDBURL).SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to MongoDB!")

	Client = client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("cluster0").Collection(collectionName)
	return collection
}

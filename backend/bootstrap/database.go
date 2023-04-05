package bootstrap

import (
	"context"
	"log"
	"time"

	"github.com/pmohanj/web-chat-app/mongo"
)

func DBinstance(MongoDBURL string) mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.NewClient(MongoDBURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("connected to MongoDB!")

	return client
}

func CloseDBinstance(client mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB is closed")
}

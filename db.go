package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	dotenv := EnvVar("CONNECT_URI")
	clientOptions := options.Client().ApplyURI(dotenv)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected")
	client.Database("goAPI").CreateCollection(context.TODO(), "Users")
	client.Database("goAPI").CreateCollection(context.TODO(), "Posts")
	return client
}

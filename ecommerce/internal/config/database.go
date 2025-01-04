package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

var (
	uri    = os.Getenv("MONGO_URI")
	dbName = os.Getenv("MONGO_DB_NAME")
)

func ConnectDatabase() {

	if uri == "" || dbName == "" {
		log.Fatalf("MONGO_URI or MONGO_DB_NAME is not set")
	}

	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	MongoClient = client
}

func GetCollection(collectionName string) *mongo.Collection {

	if uri == "" || dbName == "" {
		log.Fatal("MONGO_URI or MONGO_DB_NAME is not set")
	}

	return MongoClient.Database(dbName).Collection(collectionName)
}

package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongoDB() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	// Test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}

	MongoClient = client
	MongoDB = client.Database("db_uas")

	log.Println("MongoDB Connected")

	// Create collections and indexes
	CreateMongoCollections()
}

func CreateMongoCollections() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create achievement collection if not exists
	achievementCollection := MongoDB.Collection("achievements")

	// Create index on student_id for faster queries
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "student_id", Value: 1}},
	}

	_, err := achievementCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Println("Failed to create index on achievements:", err)
	} else {
		log.Println("MongoDB collections and indexes created successfully")
	}
}

func DisconnectMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		MongoClient.Disconnect(ctx)
		log.Println("MongoDB Disconnected")
	}
}

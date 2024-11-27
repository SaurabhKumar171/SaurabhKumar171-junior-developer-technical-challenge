package utils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// Load environment variables from .env file
func LoadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// ConnectDB connects to MongoDB using the connection string from the environment variables
func ConnectDB() *mongo.Client {
	// Load environment variables
	LoadEnvVariables()

	// Get MongoDB connection URL from the environment variable
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		log.Fatal("MongoDB URL is required in .env file")
	}

	clientOptions := options.Client().ApplyURI(mongoURL)

	// Connect to MongoDB
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")
	MongoClient = client
	return client
}

// GetCollection returns a reference to the MongoDB collection
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("rickAndMorty").Collection(collectionName)
}

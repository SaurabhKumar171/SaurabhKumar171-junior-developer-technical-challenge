package handlers

import (
	"backend/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var characterCollection *mongo.Collection

func init() {
	// Initialize the MongoDB collection
	client := utils.ConnectDB()
	characterCollection = utils.GetCollection(client, "characters") // Use global variable
}

// Struct to bind the incoming character data (matches the provided format)
type Character struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	Species  string   `json:"species"`
	Type     string   `json:"type"`
	Gender   string   `json:"gender"`
	Origin   Location `json:"origin"`
	Location Location `json:"location"`
	Image    string   `json:"image"`
	Episode  []string `json:"episode"`
	URL      string   `json:"url"`
	Created  string   `json:"created"`
}

// Location struct for nested origin and location data
type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// CreateCharacter handles the creation of a new character
func CreateCharacter(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create a new character...")

	var newCharacter Character
	err := json.NewDecoder(r.Body).Decode(&newCharacter)
	if err != nil {
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	// To Ensure mandatory fields are provided
	if newCharacter.Name == "" || newCharacter.Status == "" || newCharacter.Species == "" {
		http.Error(w, "Missing required fields: name, status, species", http.StatusBadRequest)
		log.Println("Validation failed: Missing required fields")
		return
	}

	// Insert the new character into the MongoDB collection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = characterCollection.InsertOne(ctx, newCharacter)
	if err != nil {
		http.Error(w, "Failed to insert character", http.StatusInternalServerError)
		log.Printf("Error inserting character: %v", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Character created successfully"})
}

// GetCharacters fetches all characters from the database
func GetCharacters(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching characters with pagination...")

	// Set default values for pagination
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	if page == "" {
		page = "1" // Default to page 1 if not provided
	}
	if limit == "" {
		limit = "10" // Default to 10 characters per page if not provided
	}

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}
	limitNum, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit number", http.StatusBadRequest)
		return
	}

	// Define the pagination options
	skip := (pageNum - 1) * limitNum
	options := options.Find().SetSkip(int64(skip)).SetLimit(int64(limitNum))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	totalCount, err := characterCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to count characters", http.StatusInternalServerError)
		log.Printf("Error counting characters: %v", err)
		return
	}

	cursor, err := characterCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		http.Error(w, "Failed to fetch characters", http.StatusInternalServerError)
		log.Printf("Error fetching characters: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var characters []bson.M
	if err = cursor.All(ctx, &characters); err != nil {
		http.Error(w, "Failed to parse characters", http.StatusInternalServerError)
		log.Printf("Error parsing characters: %v", err)
		return
	}

	if characters == nil {
		characters = []bson.M{}
	}

	// Return characters with pagination info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"characters": characters,
		"page":       pageNum,
		"limit":      limitNum,
		"total":      totalCount,
	})
}

// SearchCharacter handles searching for characters by name
func SearchCharacter(w http.ResponseWriter, r *http.Request) {
	log.Println("Searching characters by name...")

	// Parse pagination parameters
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	if page == "" {
		page = "1" // Default to page 1 if not provided
	}
	if limit == "" {
		limit = "10" // Default to 10 items per page if not provided
	}

	// Convert pagination parameters to integers
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		log.Println("Invalid page number:", page)
		return
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		http.Error(w, "Invalid limit number", http.StatusBadRequest)
		log.Println("Invalid limit number:", limit)
		return
	}

	// Calculate skip value for pagination
	skip := (pageNum - 1) * limitNum

	// Parse search query parameter
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		log.Println("Search request failed: 'name' parameter missing")
		return
	}

	// Perform case-insensitive search using a regular expression
	filter := bson.M{"name": bson.M{"$regex": name, "$options": "i"}}

	// Create the context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch the total number of documents matching the filter
	totalCount, err := characterCollection.CountDocuments(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to count documents", http.StatusInternalServerError)
		log.Printf("Error counting documents: %v", err)
		return
	}

	// Fetch characters with skip and limit
	cursor, err := characterCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limitNum)))
	if err != nil {
		http.Error(w, "Failed to search characters", http.StatusInternalServerError)
		log.Printf("Error searching characters: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// Parse the characters
	var characters []bson.M
	if err = cursor.All(ctx, &characters); err != nil {
		http.Error(w, "Failed to parse characters", http.StatusInternalServerError)
		log.Printf("Error parsing characters: %v", err)
		return
	}

	// Ensure characters is an empty array if no results are found
	if characters == nil {
		characters = []bson.M{}
	}

	// Return the search results along with pagination metadata
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"page":       pageNum,
		"limit":      limitNum,
		"characters": characters,
		"total":      totalCount,
	})
}

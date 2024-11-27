// routes.go
package routes

import (
	"backend/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Define the root route
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Ricky and Morty Backend!"))
	}).Methods("GET")

	// Define other routes
	r.HandleFunc("/api/getCharacters", handlers.GetCharacters).Methods("GET")
	r.HandleFunc("/api/characters", handlers.CreateCharacter).Methods("POST") // Added POST route
	r.HandleFunc("/api/characters/search", handlers.SearchCharacter).Methods("GET")

	return r
}

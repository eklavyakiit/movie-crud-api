package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string   `json:"id"`
	Isbn     string   `json:"isbn"`
	Title    string   `json:"title"`
	Director Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

// GET /movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// GET /movies/{id}
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
}

// POST /movies
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Generate a new ID that does not already exist.
	maxID := 0

	for _, item := range movies {
		id, err := strconv.Atoi(item.ID)
		if err == nil && id > maxID {
			maxID = id
		}
	}

	movie.ID = strconv.Itoa(maxID + 1)

	movies = append(movies, movie)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

// PUT /movies/{id}
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	var updatedMovie Movie

	if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for index, item := range movies {
		if item.ID == params["id"] {
			// Keep the original ID from the URL.
			updatedMovie.ID = item.ID
			movies[index] = updatedMovie

			json.NewEncoder(w).Encode(updatedMovie)
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
}

// DELETE /movies/{id}
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)

			json.NewEncoder(w).Encode(movies)
			return
		}
	}

	http.Error(w, "Movie not found", http.StatusNotFound)
}

// UPDATE MOVIES - delete first and then overwrite

func updateMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the movie ID from the URL
	params := mux.Vars(r)
	id := params["id"]

	// Store the updated data sent by the user
	var updatedMovie Movie

	err := json.NewDecoder(r.Body).Decode(&updatedMovie)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Find the existing movie and replace its data
	for index, item := range movies {
		if item.ID == id {
			updatedMovie.ID = item.ID
			movies[index] = updatedMovie

			json.NewEncoder(w).Encode(updatedMovie)
			return
		}
	}

	// If no movie with that ID exists
	http.Error(w, "Movie not found", http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()

	// Initial movie data
	movies = append(movies,
		Movie{
			ID:    "1",
			Isbn:  "100001",
			Title: "3 Idiots",
			Director: Director{
				FirstName: "Rajkumar",
				LastName:  "Hirani",
			},
		},
		Movie{
			ID:    "2",
			Isbn:  "100002",
			Title: "Dangal",
			Director: Director{
				FirstName: "Nitesh",
				LastName:  "Tiwari",
			},
		},
		Movie{
			ID:    "3",
			Isbn:  "100003",
			Title: "Shershaah",
			Director: Director{
				FirstName: "Vishnuvardhan",
				LastName:  "",
			},
		},
		Movie{
			ID:    "4",
			Isbn:  "100004",
			Title: "Zindagi Na Milegi Dobara",
			Director: Director{
				FirstName: "Zoya",
				LastName:  "Akhtar",
			},
		},
		Movie{
			ID:    "5",
			Isbn:  "100005",
			Title: "Lagaan",
			Director: Director{
				FirstName: "Ashutosh",
				LastName:  "Gowariker",
			},
		},
		Movie{
			ID:    "6",
			Isbn:  "100006",
			Title: "Taare Zameen Par",
			Director: Director{
				FirstName: "Aamir",
				LastName:  "Khan",
			},
		},
	)

	// API routes
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
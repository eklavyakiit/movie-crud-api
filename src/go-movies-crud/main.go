package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct{
	ID string `json: "id"`
	Isbn string `json: "isbn"`
	Title string `json: "title"`
	Director string `json: "director"`
}

type Director struct{
	FirstName string `json: "firstname"`
	LastName string `json: "lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies{
		if item.ID == params["id"]{
			movies = append(movies[:index], movies[index+1]...)
			break
		}
	}
}

func getMovie(w htt.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	
}
func main(){
	r := mux.NewRouter()


	movies = append(movies, Movie{ID: "1", Isbn: "100001", Title: "3 Idiots", Director: &Director{Firstname: "Rajkumar", Lastname: "Hirani"}})
	movies = append(movies, Movie{ID: "2", Isbn: "100002", Title: "Dangal", Director: &Director{Firstname: "Nitesh", Lastname: "Tiwari"}})
	movies = append(movies, Movie{ID: "3", Isbn: "100003", Title: "Shershaah", Director: &Director{Firstname: "Vishnuvardhan", Lastname: ""}})
	movies = append(movies, Movie{ID: "4", Isbn: "100004", Title: "Zindagi Na Milegi Dobara", Director: &Director{Firstname: "Zoya", Lastname: "Akhtar"}})
	movies = append(movies, Movie{ID: "5", Isbn: "100005", Title: "Lagaan", Director: &Director{Firstname: "Ashutosh", Lastname: "Gowariker"}})
	movies = append(movies, Movie{ID: "6", Isbn: "100006", Title: "Taare Zameen Par", Director: &Director{Firstname: "Aamir", Lastname: "Khan"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovies).Methods("POST")
	r.HandleFunc("/movies{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting Server at port 8080")
	//func http.ListenAndServe(addr string, handler http.Handler) error
	//log.fatal means that is is an essential thing if server is not listening its a fatal error
	log.Fatal(http.ListenAndServe(":8080", r))


}
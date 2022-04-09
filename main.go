package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
    DB_USER     = "postgres"
    DB_PASSWORD = "root"
    DB_NAME     = "postgres"
)

// DB set up
func setupDB() *sql.DB {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
    db, err := sql.Open("postgres", dbinfo)

    checkErr(err)

    return db
}


type Movie struct{
	MovieID string `json:"movieid"`
	MovieName string `json:"moviename"`
}

type JsonResponse struct {
    Type    string `json:"type"`
    Data    []Movie `json:"data"`
    Message string `json:"message"`
}


func main(){
	router := mux.NewRouter();

	router.HandleFunc("/movies",GetMovies).Methods("GET")
	router.HandleFunc("/movies/{movieid}",GetMovie).Methods("GET")
	router.HandleFunc("/movies",CreateMovie).Methods("POST")
	router.HandleFunc("/movies/{movieid}",DeleteMovie).Methods("DELETE")
	router.HandleFunc("/movies", DeleteAllMovies).Methods("DELETE")

	fmt.Println("Server running on port 6000")
	log.Fatal(http.ListenAndServe(":6000",router))
}

// Function for handling messages
func printMessage(message string) {
    fmt.Println("")
    fmt.Println(message)
    fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

// Get all movies

// response and request handlers
func GetMovies(w http.ResponseWriter, r *http.Request) {
    db := setupDB()

    printMessage("Getting movies...")

    // Get all movies from movies table that don't have movieID = "1"
    rows, err := db.Query("SELECT * FROM movies")

    // check errors
    checkErr(err)

    // var response []JsonResponse
    var movies []Movie

    // Foreach movie
    for rows.Next() {
        var id int
        var movieID string
        var movieName string

        err = rows.Scan(&id, &movieID, &movieName)

        // check errors
        checkErr(err)

        movies = append(movies, Movie{MovieID: movieID, MovieName: movieName})
    }

    var response = JsonResponse{Type: "success", Data: movies}

    json.NewEncoder(w).Encode(response)
}

// get specific movie
func GetMovie(w http.ResponseWriter, r *http.Request){

	params := mux.Vars(r)
    movieID := params["movieid"]
	println(movieID)
	db := setupDB()

    printMessage("Getting Specific Movie ...")
	
    // Get specific movie from movies table 
    rows, err := db.Query("SELECT DISTINCT * FROM movies WHERE movieID = $1", movieID)
	// println(rows)
	checkErr(err)

	    // var response []JsonResponse
    var movies []Movie

    // Foreach movie
    for rows.Next() {
        var id int
        var movieID string
        var movieName string

        err = rows.Scan(&id, &movieID, &movieName)

        // check errors
        checkErr(err)

		movies = append(movies, Movie{MovieID: movieID, MovieName: movieName})
	}

	var response = JsonResponse{Type: "success", Data: movies}

    json.NewEncoder(w).Encode(response)

}

// Create a movie

// response and request handlers
func CreateMovie(w http.ResponseWriter, r *http.Request) {
    movieID := r.FormValue("movieid")
    movieName := r.FormValue("moviename")

    var response = JsonResponse{}

    if movieID == "" || movieName == "" {
        response = JsonResponse{Type: "error", Message: "You are missing movieID or movieName parameter."}
    } else {
        db := setupDB()

        printMessage("Inserting movie into DB")

        fmt.Println("Inserting new movie with ID: " + movieID + " and name: " + movieName)

        var lastInsertID int
    err := db.QueryRow("INSERT INTO movies(movieID, movieName) VALUES($1, $2) returning id;", movieID, movieName).Scan(&lastInsertID)

    // check errors
    checkErr(err)

    response = JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}
    }

    json.NewEncoder(w).Encode(response)
}

// Delete a movie

// response and request handlers
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    movieID := params["movieid"]

    var response = JsonResponse{}

    if movieID == "" {
        response = JsonResponse{Type: "error", Message: "You are missing movieID parameter."}
    } else {
        db := setupDB()

        printMessage("Deleting movie from DB")

        _, err := db.Exec("DELETE FROM movies where movieID = $1", movieID)

        // check errors
        checkErr(err)

        response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
    }

    json.NewEncoder(w).Encode(response)
}


// Delete all movies

// response and request handlers
func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
    db := setupDB()

    printMessage("Deleting all movies...")

    _, err := db.Exec("DELETE FROM movies")

    // check errors
    checkErr(err)

    printMessage("All movies have been deleted successfully!")

    var response = JsonResponse{Type: "success", Message: "All movies have been deleted successfully!"}

    json.NewEncoder(w).Encode(response)
}
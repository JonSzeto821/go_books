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

// create type schemas
type Book struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Author       string        `json:"author"`
	Isbn         string        `json:"isbn"`
	Year         int           `json:"year"`
	PersonalCopy *PersonalCopy `json:"personalCopy"`
}

type PersonalCopy struct {
	ReadStatus    bool   `json:"readStatus"`
	ContentFormat string `json:"contentFormat"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	found := false

	for _, book := range books {
		if book.ID == params["id"] {
			found = true
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	if !found {
		handleError(w, http.StatusNotFound, "No results found with ID")
	}
}

// TODO move to a utils package
func handleError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.WriteHeader(statusCode)
	errorResponse := ErrorResponse{Message: errorMessage}
	json.NewEncoder(w).Encode(errorResponse)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, book := range books {
		if book.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = generateID()
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	found := false

	// approach -> delete the original record, append new record at the end of slice
	for index, book := range books {
		if book.ID == params["id"] {
			found = true
			books = append(books[:index], books[index+1:]...)

			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = generateID()
			books = append(books, book)
			json.NewEncoder(w).Encode(books)
			return
		}
	}

	if !found {
		handleError(w, http.StatusNotFound, "No results found with ID")
	}
}

// TODO move to a utils package
func generateID() string {
	return strconv.Itoa(rand.Intn(1000000000))
}

func main() {
	//instantiate a new router
	r := mux.NewRouter()

	// create sample book records //TODO move sample data to standalone file and import back in
	books = append(books, Book{ID: generateID(), Title: "Dune", Author: "Frank Herbert", Isbn: "9780399128967", Year: 1965, PersonalCopy: &PersonalCopy{ReadStatus: true, ContentFormat: "Kindle"}})

	books = append(books, Book{ID: generateID(), Title: "The Three-Body Problem", Author: "Liu Cixin", Isbn: "9780765377067", Year: 2014, PersonalCopy: &PersonalCopy{ReadStatus: false, ContentFormat: "Kindle"}})

	books = append(books, Book{ID: generateID(), Title: "The Name of the Wind", Author: "Patrick Rothfuss", Isbn: "9780756404079", Year: 2008, PersonalCopy: &PersonalCopy{ReadStatus: true, ContentFormat: "Paperback"}})

	// GET ALL
	r.HandleFunc("/books", getBooks).Methods("GET")

	// GET BY ID
	r.HandleFunc("/books/{id}", getBook).Methods("GET")

	// CREATE
	r.HandleFunc("/books", createBook).Methods("POST")

	// UPDATE
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")

	// DELETE
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}

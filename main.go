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

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, book := range books {
		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
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
	book.ID = strconv.Itoa(rand.Intn(1000000000))
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// approach -> delete the original record, append new record at the end of slice
	for index, book := range books {
		if book.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)

			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			books = append(books, book)
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	//instantiate a new router
	r := mux.NewRouter()

	// create sample book records
	books = append(books, Book{ID: "1", Title: "Dune", Author: "Frank Herbert", Isbn: "9780399128967", Year: 1965, PersonalCopy: &PersonalCopy{ReadStatus: true, ContentFormat: "Kindle"}})

	books = append(books, Book{ID: "2", Title: "The Three-Body Problem", Author: "Liu Cixin", Isbn: "9780765377067", Year: 2014, PersonalCopy: &PersonalCopy{ReadStatus: false, ContentFormat: "Kindle"}})

	books = append(books, Book{ID: "3", Title: "The Name of the Wind", Author: "Patrick Rothfuss", Isbn: "9780756404079", Year: 2008, PersonalCopy: &PersonalCopy{ReadStatus: true, ContentFormat: "Paperback"}})

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

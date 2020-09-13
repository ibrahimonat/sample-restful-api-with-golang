package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type book struct {
	ID    uint64  `json:"ID"`
	Name  string  `json:"Name"`
	Price float32 `json:"Price"`
}

// Sample Book Database
type allBooks []book

var books allBooks

func seedBooks() {
	books = allBooks{
		{
			ID:    1,
			Name:  "Introducing Go: Build Reliable, Scalable Programs",
			Price: 20.69,
		},
		{
			ID:    2,
			Name:  "The Hitchhiker's Guide to the Galaxy",
			Price: 40.2,
		},
		{
			ID:    3,
			Name:  "1984",
			Price: 30.55,
		},
	}
}

// BOOK API

func homePage(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Welcome API Documentation!")
}

// POST
func insert(res http.ResponseWriter, req *http.Request) {
	var newBook book
	reqBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintf(res, "Please enter data with the book name and price only in order to create!")
	}

	json.Unmarshal(reqBody, &newBook)
	books = append(books, newBook)
	res.WriteHeader(http.StatusCreated)

	json.NewEncoder(res).Encode(newBook)
}

// GET
func get(res http.ResponseWriter, req *http.Request) {
	bookID, err := strconv.ParseUint(mux.Vars(req)["Id"], 10, 64)

	if err != nil {
		fmt.Fprintf(res, "Please enter data with the book ID only in order to get!")
	}

	for _, book := range books {
		if book.ID == bookID {
			json.NewEncoder(res).Encode(book)
		}
	}
}

// GET
func getAll(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(books)
}

// PATCH
func update(res http.ResponseWriter, req *http.Request) {
	bookID, err := strconv.ParseUint(mux.Vars(req)["Id"], 10, 64)

	if err != nil {
		fmt.Fprintf(res, "Please enter data with the book ID only in order to update!")
	}

	var updatedBook book
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(res, "Please enter data with the book name and price only in order to update!")
	}

	json.Unmarshal(reqBody, &updatedBook)

	for i, book := range books {
		if book.ID == bookID {
			book.Name = updatedBook.Name
			book.Price = updatedBook.Price
			books = append(books[:i], book)
			json.NewEncoder(res).Encode(book)
		}
	}
}

// DELETE
func delete(res http.ResponseWriter, req *http.Request) {
	bookID, err := strconv.ParseUint(mux.Vars(req)["Id"], 10, 64)

	if err != nil {
		fmt.Fprintf(res, "Please enter data with the book ID only in order to delete!")
	}

	for i, book := range books {
		if book.ID == bookID {
			books = append(books[:i], books[i+1:]...)
			fmt.Fprintf(res, "The book with Id %v has been deleted successfully!", bookID)
		}
	}
}

func main() {
	seedBooks()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/book", insert).Methods("POST")
	router.HandleFunc("/books", getAll).Methods("GET")
	router.HandleFunc("/books/{Id}", get).Methods("GET")
	router.HandleFunc("/books/{Id}", update).Methods("PATCH")
	router.HandleFunc("/books/{Id}", delete).Methods("DELETE")
	fmt.Print("Book API is running...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

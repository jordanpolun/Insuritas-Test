// Every Go program needs this, I guess
package main

// Import the libraries we need
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// This is the struct of our JSON object
type Book struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	NumPages int       `json:"num_pages"`
	PubDate  time.Time `json:"pub_date"`
}

// Initializing our array of JSON objects
// Traditionally, global variables are a bad idea, but I looked it up
// and people say that it's okay in Go
var library []Book

// This function will determine what happens when the homepage endpoint is hit
func HomePage(w http.ResponseWriter, r *http.Request) {
	// This writes to the page so we know we've hit the homepage
	fmt.Fprintf(w, "Welcome to the HomePage!")

	// This tells us in the terminal that someone has hit the homepage
	fmt.Println("Endpoint Hit: homePage")
}

// This function will route function-calls based on the endpoint hit
func HandleRequests() {

	// The base URL will run the homepage() function
	http.HandleFunc("/", HomePage)

	// The /library could mean multiple things, send the requests to a function
	http.HandleFunc("/library", HandleLibraryRequests)

	// This starts the server on localhost:10000 and returns an error if it can't
	log.Fatal(http.ListenAndServe(":10000", nil))
}

// This function will determine which type of request we've received and do what it needs to do
func HandleLibraryRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Serve the resource.
		ReturnLibrary(w, r)
	case http.MethodPost:
		// Create a new record.
		AddBook(w, r)
	case http.MethodPut:
		// Update an existing record.
		EditBook(w, r)
	case http.MethodDelete:
		// Remove the record.
		RemoveBook(w, r)
	default:
		// Give an error message.
		fmt.Fprintf(w, "Only GET, POST, PUT, and DELETE requests are supported")
	}
}

func ReturnLibrary(w http.ResponseWriter, r *http.Request) {

	// Letting us know that something went to the URL that caused this function to run
	fmt.Println("Endpoint Hit: returnLibrary")

	// This is the encoder we will use to write our JSON object
	json_encoder := json.NewEncoder(w)

	// Find what id was specified in query, if any
	index, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || index >= len(library) {
		// If any error, just show the whole library, print error to console
		json_encoder.Encode(library)
		fmt.Println(err)
	} else {
		// Show the book specified
		json_encoder.Encode(library[index])
	}
}

func AddBook(w http.ResponseWriter, r *http.Request) {

	// Create a new book variable and decode from the request body
	var b Book
	err := json.NewDecoder(r.Body).Decode(&b)

	// If there was no error, add that book to the library
	if err != nil {
		fmt.Fprintf(w, "Failed to POST")
		fmt.Println(err)
	} else {
		// Change b Id to auto increment from most recent entry
		b.Id = library[len(library)-1].Id + 1
		library = append(library, b)

		fmt.Fprintf(w, "Successful POST")
		fmt.Println(b.Title, "has been added to the library")
	}
}

func EditBook(w http.ResponseWriter, r *http.Request) {

	// Create a new book variable and decode from request body
	var b Book
	err := json.NewDecoder(r.Body).Decode(&b)

	// If there was no error, place that book in the library at its index
	if err != nil {
		fmt.Fprintf(w, "Failed to PUT")
		fmt.Println(err)
	} else {
		// Replace book at current index with this new (edited) book
		library[b.Id] = b

		fmt.Fprintf(w, "Successful PUT")
		fmt.Println(b.Title, "has been edited in the library")
	}

}

func RemoveBook(w http.ResponseWriter, r *http.Request) {
	// Find what id was specified in query, if any
	index, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || index >= len(library) {
		// If any error, just say it failed, print error to console
		fmt.Fprintf(w, "Failed to DELETE index", r.FormValue("id"))
		fmt.Println(err)
	} else {
		// Delete the book specified
		library[index] = library[len(library)-1] // Copy last book to this index
		library = library[:len(library)-1]       // Remove last index of library
		library[index].Id = index                // Reset id of element that was moved
    fmt.Fprintf(w, "Deleted index", r.FormValue("id"))
	}
}

func main() {

	// Declare what is in our array of JSON objects
	library = []Book{
		{Id: 0, Title: "Charlotte's Web", Author: "E. B. White", NumPages: 192, PubDate: time.Date(1952, time.October, 15, 0, 0, 0, 0, time.UTC)},
		{Id: 1, Title: "The Very Hungry Caterpillar", Author: "Eric Carle", NumPages: 22, PubDate: time.Date(1969, time.June, 3, 0, 0, 0, 0, time.UTC)},
		{Id: 2, Title: "Charlie and the Chocolate Factory", Author: "Roald Dahl", NumPages: 178, PubDate: time.Date(1964, time.January, 17, 0, 0, 0, 0, time.UTC)},
		{Id: 3, Title: "The Phantom Tollbooth", Author: "Norton Juster", NumPages: 255, PubDate: time.Date(1961, time.October, 12, 0, 0, 0, 0, time.UTC)},
	}

	// This method does the rerouting and heavy lifting
	HandleRequests()
}

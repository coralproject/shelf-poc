package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func main() {

	// Check connection to Mongo.
	log.Println("Checking connection to MongoDB")
	mgoDB, err := mongoConnection("coral_test")
	if err != nil {
		err = errors.Wrap(err, "Could not connect to Mongo")
		log.Fatal(err)
	}
	mgoDB.CloseMGO("Mongo")

	// ListenAndServe starts an HTTP server with a given address and
	// handler defined in NewRouter
	log.Println("Starting JSON API, listening on port 8080...")
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))

}

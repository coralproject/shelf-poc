package main

import (
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

var (
	shelfHost   string
	numRequests int
	numDocs     int
	results     Results
	timeOut     int
)

func init() {

	// Import the number of requests we will execute and analyze.
	numReqEnv, present := os.LookupEnv("STATS_NUM_REQ")
	if !present {
		log.Fatal("The STATS_NUM_REQ environmental var. is not defined")
	}

	// Convert the number of requests to an integer.
	var err error
	numRequests, err = strconv.Atoi(numReqEnv)
	if err != nil {
		err := errors.Wrap(err, "Could not parse the number of requests")
		log.Fatal(err)
	}

	// Import the shelf host.
	shelfHost, present = os.LookupEnv("STATS_SHELF_HOST")
	if !present {
		log.Fatal("The STATS_SHELF_HOST environmental var. is not defined")
	}

	// Import the number of documents we will analyze in Mongo.
	numDocEnv, present := os.LookupEnv("STATS_NUM_DOCS")
	if !present {
		log.Fatal("The STATS_NUM_DOCS environmental var. is not defined")
	}

	// Convert the number of documents to an integer.
	numDocs, err = strconv.Atoi(numDocEnv)
	if err != nil {
		err := errors.Wrap(err, "Could not parse the number of documents")
		log.Fatal(err)
	}

}

func main() {

	log.Println("Starting to gather shelf stats")
	results = Results{}
	timeOut = 0

	log.Println("Executing user comment view requests")
	if err := userComments(); err != nil {
		log.Fatal(err)
	}

	log.Println("Executing user asset view requests")
	if err := userAssets(); err != nil {
		log.Fatal(err)
	}

	log.Println("Executing asset view requests")
	if err := assetCommentsAuthors(); err != nil {
		log.Fatal(err)
	}

	log.Println("Executing parent view requests")
	if err := parentedComments(); err != nil {
		log.Fatal(err)
	}

	log.Println("Executing grandparent view requests")
	if err := grandParentedComments(); err != nil {
		log.Fatal(err)
	}

	log.Println("Executing great-grandparent view requests")
	if err := grGrandParentedComments(); err != nil {
		log.Fatal(err)
	}

	// Output the results.
	printResults()
}

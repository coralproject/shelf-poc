package main

import (
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

var (
	numDoc     int
	spongeHost string
	jobs       chan Job
	results    chan error
)

const (
	numWorkers = 100
	numJobs    = 200
)

func init() {

	// Import the number of number of random documents we will generate and anlyze
	numDocEnv, present := os.LookupEnv("SHELF_NUM_DOC")
	if !present {
		log.Fatal("The SHELF_NUM_DOC environmental var. is not defined")
	}

	// Convert the number of documents to an integer
	var err error
	numDoc, err = strconv.Atoi(numDocEnv)
	if err != nil {
		err := errors.Wrap(err, "Could not parse the number of documents")
		log.Fatal(err)
	}

	// Import the sponge host via environmental var.
	spongeHost, present = os.LookupEnv("SHELF_SPONGE_HOST")
	if !present {
		log.Fatal("The SHELF_SPONGE_HOST environmental var. is not defined")
	}

}

func main() {

	jobs = make(chan Job, numJobs)
	results = make(chan error, numJobs)
	handleErrors()

	log.Println("Get numbers of dummy items to be generated")
	numComments, numUsers, numAssets := generateItemNumbers()

	log.Println("Start workers")
	for w := 1; w <= numWorkers; w++ {
		go worker(jobs, results)
	}

	log.Println("Generate dummy users")
	if err := generateUsers(numUsers); err != nil {
		log.Fatal(err)
	}

	log.Println("Generate dummy assets")
	if err := generateAssets(numAssets); err != nil {
		log.Fatal(err)
	}

	log.Println("Generate dummy comments")
	if err := generateComments(numComments, numUsers, numAssets); err != nil {
		log.Fatal(err)
	}

	close(jobs)
}

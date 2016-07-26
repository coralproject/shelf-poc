package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/pkg/errors"
)

var (
	numDoc     int
	mongoHost  string
	jobs       chan Job
	results    chan error
	numWorkers int
	numJobs    int
)

func init() {

	// Import the number of number of random documents we will generate and anlyze.
	numDocEnv, present := os.LookupEnv("SHELF_NUM_DOC")
	if !present {
		log.Fatal("The SHELF_NUM_DOC environmental var. is not defined")
	}

	// Convert the number of documents to an integer.
	var err error
	numDoc, err = strconv.Atoi(numDocEnv)
	if err != nil {
		err := errors.Wrap(err, "Could not parse the number of documents")
		log.Fatal(err)
	}

	// Import the sponge host.
	mongoHost, present = os.LookupEnv("SHELF_MONGO_HOST")
	if !present {
		log.Fatal("The SHELF_MONGO_HOST environmental var. is not defined")
	}

	// Import the number of number of workers and jobs.
	numWorkersEnv, present := os.LookupEnv("SHELF_WORKERS")
	if !present {
		log.Fatal("The SHELF_WORKERS environmental var. is not defined")
	}
	numWorkers, err = strconv.Atoi(numWorkersEnv)
	if err != nil {
		err := errors.Wrap(err, "Could not parse the number of workers")
		log.Fatal(err)
	}
	numJobsEnv, present := os.LookupEnv("SHELF_JOBS")
	if !present {
		log.Fatal("The SHELF_JOBS environmental var. is not defined")
	}
	numJobs, err = strconv.Atoi(numJobsEnv)
	if err != nil {
		err := errors.Wrap(err, "Could not parse the number of jobs")
		log.Fatal(err)
	}
}

func main() {

	// Make the channels for handling data imports.
	jobs = make(chan Job, numJobs)
	results = make(chan error, numJobs)
	handleErrors()

	// Check connection to Mongo.
	log.Println("Checking connection to MongoDB")
	mgoDB, err := mongoConnection("coral_test")
	if err != nil {
		err = errors.Wrap(err, "Could not connect to Mongo")
		log.Fatal(err)
	}
	mgoDB.CloseMGO("Mongo")

	// Initialize the graph database.
	mongoHostPort := fmt.Sprintf("%s:27017", mongoHost)
	if err := graph.InitQuadStore("mongo", mongoHostPort, nil); err != nil {
		err = errors.Wrap(err, "Could not initialize quad store")
		log.Fatal(err)
	}

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

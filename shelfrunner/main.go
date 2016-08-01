package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
)

var (
	fills     chan FillRel
	payload   string
	mongoHost string
)

func init() {

	// Import the number of number of random documents we will generate and anlyze.
	var present bool
	payload, present = os.LookupEnv("PAYLOAD")
	if !present {
		log.Fatal("The PAYLOAD environmental var. is not defined")
	}

}

func main() {

	// Parse the asset ID.
	var ironInput IronInput
	err := json.Unmarshal([]byte(payload), &ironInput)
	mongoHost = ironInput.MongoHost

	// Check connection to Mongo.
	log.Println("Checking connection to MongoDB")
	mgoDB, err := mongoConnection("coral_test")
	if err != nil {
		err = errors.Wrap(err, "Could not connect to Mongo")
		log.Fatal(err)
	}
	mgoDB.CloseMGO("Mongo")

	// Start listening for relationships that need to be embedded in
	// saved views.
	fills = make(chan FillRel, 1000)
	fillRelationships()

	// Generate and save a view.
	itemIDs, fillRels, err := getItemsOnAsset(ironInput.Asset)
	if err != nil {
		log.Fatal(err)
	}
	output, err := saveView(itemIDs, fillRels)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Saved view as: %s\n", output.Collection)

	if len(fills) > 0 {
		time.Sleep(1 * time.Millisecond)
	}
}

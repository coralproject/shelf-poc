package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/pkg/errors"
)

// sendToMongo stores an item in Mongo.
func sendToMongo(item Item) error {

	// Create a session copy.
	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

	// Insert or update the item.
	f := func(c *mgo.Collection) error {
		q := bson.M{"_id": item.ID}
		_, err := c.Upsert(q, item)
		return err
	}

	if err := mgoDB.ExecuteMGO("Mongo", "coral_items", f); err != nil {
		return err
	}

	return nil
}

// sendToCayley imports relationships into Cayley.
func sendToCayley(tx *graph.Transaction) error {

	// Get the connection to Cayley.
	mongoHostPort := fmt.Sprintf("%s:27017", mongoHost)
	store, err := cayley.NewGraph("mongo", mongoHostPort, nil)
	if err != nil {
		return errors.Wrap(err, "Could not open and use the Cayley DB")
	}
	defer store.Close()

	// Apply all the transactions that we have accumulated.
	if err := store.ApplyTransaction(tx); err != nil {
		return errors.Wrap(err, "Could not execute transaction")
	}
	return nil
}

// worker processes a job on the jobs channel and publishes the result to the
// results channel.
func worker(jobsIn <-chan Job, resultsIn chan<- error) {
	for j := range jobsIn {
		if err := sendToMongo(j.Data); err != nil {
			resultsIn <- err
		}
		if err := sendToCayley(j.Tx); err != nil {
			resultsIn <- err
		}
	}
}

// handleErrors listens for errors generated in uploading items.
func handleErrors() {
	go func() {
		for err := range results {
			if err != nil {
				panic(err)
			}
		}
	}()
}

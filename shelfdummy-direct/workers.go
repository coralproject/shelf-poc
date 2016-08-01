package main

import (
	"math/rand"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/pkg/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

// worker processes a job on the jobs channel and publishes the result to the
// results channel.
func worker(jobsIn <-chan Job, resultsIn chan<- error) {
	for j := range jobsIn {
		if err := sendToMongo(j.Data); err != nil {
			resultsIn <- err
		}
		txMutex.Lock()
		for _, quad := range j.Quads {
			tx.AddQuad(quad)
		}
		txMutex.Unlock()
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

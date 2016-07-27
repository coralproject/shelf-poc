package main

import (
	"math/rand"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ardanlabs/kit/db"
	"github.com/cayleygraph/cayley"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/cayleygraph/cayley/quad"
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

// parentQuad returns a parent comment randomly.
func parentQuad(job Job) (quad.Quad, error) {

	// Get parent relationship if necessary.
	if rand.Intn(2) == 1 && job.Type == "comment" {
		docMutex.Lock()
		commentCount := docCount["comment"]
		docMutex.Unlock()
		if commentCount > 2 {
			mgoDB, err := db.NewMGO("Mongo", "test")
			if err != nil {
				return quad.Quad{}, errors.Wrap(err, "Could not connect to MongoDB")
			}
			parent, err := retrieveRandComment(commentCount, mgoDB)
			if err != nil {
				results <- errors.Wrap(err, "Could not retrieve rand. comment")
			}
			mgoDB.CloseMGO("Mongo")
			return cayley.Quad(job.Data.ID.Hex(), "parented_by", parent.ID.Hex(), ""), nil
		}
	}

	return quad.Quad{}, nil
}

// worker processes a job on the jobs channel and publishes the result to the
// results channel.
func worker(jobsIn <-chan Job, resultsIn chan<- error) {
	for j := range jobsIn {
		if err := sendToMongo(j.Data); err != nil {
			resultsIn <- err
		}
		additionalQuad, err := parentQuad(j)
		if err != nil {
			resultsIn <- err
		}
		if additionalQuad != (quad.Quad{}) {
			txMutex.Lock()
			tx.AddQuad(additionalQuad)
			txMutex.Unlock()
		}
		docMutex.Lock()
		docCount[j.Type]++
		docMutex.Unlock()
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

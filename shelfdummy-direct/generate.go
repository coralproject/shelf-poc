package main

import (
	"math/rand"
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/cayleygraph/cayley"
	"github.com/pkg/errors"

	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// letterRunes includes letters from which the random names/content will be generated.
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generateItemNumbers generates numbers of users, comments, and assets
// that we will generate for the dummy data.  Based on NYT proportions,
// we will generate 80% comments, 15% users, and 5% assets within the
// whole of the Item collection
func generateItemNumbers() (int, int, int) {
	numComments := (numDoc * 4) / 5
	numUsers := (numDoc * 3) / 20
	numAssets := numDoc / 20
	return numComments, numUsers, numAssets
}

// generateUsers generates a number, num, of random items of type coral_user.
func generateUsers(num int) error {
	currentNum := 1
	for currentNum <= num {

		// Generate generic data.
		data := ItemData{
			ID:   currentNum,
			Name: RandStringRunes(10),
		}
		item := Item{
			ID:      bson.NewObjectId(),
			Type:    "coral_user",
			Version: 1,
			Data:    data,
		}

		// Generate relationships.
		t := cayley.NewTransaction()
		t.AddQuad(cayley.Quad(item.ID.Hex(), "is_type", "coral_user", ""))

		// Send the job to the workers.
		job := Job{
			Data: item,
			Type: "user",
			Tx:   t,
		}
		jobs <- job
		currentNum++
	}
	return nil
}

// generateAssets generates a number, num, of random items of type coral_asset.
func generateAssets(num int) error {
	currentNum := 1
	for currentNum <= num {

		// Generate generic item data.
		data := ItemData{
			ID:   currentNum,
			Name: RandStringRunes(10),
		}
		item := Item{
			ID:      bson.NewObjectId(),
			Type:    "coral_asset",
			Version: 1,
			Data:    data,
		}

		// Generate relationships.
		t := cayley.NewTransaction()
		t.AddQuad(cayley.Quad(item.ID.Hex(), "is_type", "coral_asset", ""))

		// Send the job to the workers.
		job := Job{
			Data: item,
			Type: "asset",
			Tx:   t,
		}
		jobs <- job
		currentNum++
	}
	return nil
}

// generateComments generates a number, num, of random items of type coral_comment.
func generateComments(numComments, numUsers, numAssets int) error {
	currentNum := 1
	for currentNum <= numComments {

		// Generate item specific data.
		data := ItemData{
			ID:   currentNum,
			Name: RandStringRunes(10),
			Body: RandStringRunes(20),
		}
		item := Item{
			ID:      bson.NewObjectId(),
			Type:    "coral_comment",
			Version: 1,
			Data:    data,
		}

		// Generate relationships.
		t := cayley.NewTransaction()
		mgoDB, err := db.NewMGO("Mongo", "test")
		if err != nil {
			return errors.Wrap(err, "Could not connect to MongoDB")
		}
		t.AddQuad(cayley.Quad(item.ID.Hex(), "is_type", "coral_comment", ""))

		// Add contextualized with relationship.
		asset, err := retrieveRandAsset(numDoc/20, mgoDB)
		if err != nil {
			results <- errors.Wrap(err, "Could not retrieve rand. asset")
		}
		t.AddQuad(cayley.Quad(item.ID.Hex(), "contextualized_with", asset.ID.Hex(), ""))

		// Add authored by relationship.
		author, err := retrieveRandUser((numDoc*3)/20, mgoDB)
		if err != nil {
			results <- errors.Wrap(err, "Could not retrieve rand. author")
		}
		t.AddQuad(cayley.Quad(item.ID.Hex(), "authored_by", author.ID.Hex(), ""))

		// Add parented by relationship.
		if rand.Intn(2) == 1 {
			parent, err := retrieveRandComment(currentNum, mgoDB)
			if err != nil {
				results <- errors.Wrap(err, "Could not retrieve rand. comment")
			}
			t.AddQuad(cayley.Quad(item.ID.Hex(), "parented_by", parent.ID.Hex(), ""))
		}
		mgoDB.CloseMGO("Mongo")

		// Send the job to the workers.
		job := Job{
			Data: item,
			Type: "comment",
			Tx:   t,
		}
		jobs <- job
		currentNum++
	}
	return nil
}

// RandStringRunes generates a random string for our dummy data.
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

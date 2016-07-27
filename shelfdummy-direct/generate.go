package main

import (
	"math/rand"
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"

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
		var quads []quad.Quad
		quads = append(quads, cayley.Quad(item.ID.Hex(), "is_type", "coral_user", ""))

		// Keep the object ID in memory.
		userIDs = append(userIDs, item.ID)

		// Send the job to the workers.
		job := Job{
			Data:  item,
			Type:  "user",
			Quads: quads,
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

		// Keep the object ID in memory.
		assetIDs = append(assetIDs, item.ID)

		// Generate relationships.
		var quads []quad.Quad
		quads = append(quads, cayley.Quad(item.ID.Hex(), "is_type", "coral_asset", ""))

		// Send the job to the workers.
		job := Job{
			Data:  item,
			Type:  "asset",
			Quads: quads,
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

		// Keep the comment ID in memory, in some cases (randomly).
		switch {
		case len(commentIDs) < 100:
			commentIDs = append(commentIDs, item.ID)
		case len(commentIDs) >= 100:
			commentIDs[rand.Intn(100)] = item.ID
		}

		// Generate relationships.
		var quads []quad.Quad
		quads = append(quads, cayley.Quad(item.ID.Hex(), "is_type", "coral_comment", ""))

		// Add contextualized with relationship.
		assetID := assetIDs[rand.Intn(numDoc/20)]
		quads = append(quads, cayley.Quad(item.ID.Hex(), "contextualized_with", assetID.Hex(), ""))

		// Add authored by relationship.
		authorID := userIDs[rand.Intn((numDoc*3)/20)]
		quads = append(quads, cayley.Quad(item.ID.Hex(), "authored_by", authorID.Hex(), ""))

		// Get parent relationship if necessary.
		if rand.Intn(2) == 1 {
			if len(commentIDs) > 2 {
				commentID := commentIDs[rand.Intn(len(commentIDs))]
				quads = append(quads, cayley.Quad(item.ID.Hex(), "parented_by", commentID.Hex(), ""))
			}
		}

		// Send the job to the workers.
		job := Job{
			Data:  item,
			Type:  "comment",
			Quads: quads,
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

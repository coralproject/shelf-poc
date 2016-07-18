package main

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
		user := User{
			ID:   currentNum,
			Name: RandStringRunes(10),
		}
		payload, err := json.Marshal(user)
		if err != nil {
			return errors.Wrap(err, "Could not marshal JSON user")
		}
		job := Job{
			Data: payload,
			Type: "user",
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
		asset := Asset{
			ID:    currentNum,
			Title: RandStringRunes(20),
		}
		payload, err := json.Marshal(asset)
		if err != nil {
			return errors.Wrap(err, "Could not marshal JSON asset")
		}
		job := Job{
			Data: payload,
			Type: "asset",
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
		comment := Comment{
			ID:      currentNum,
			UserID:  rand.Intn(numUsers),
			AssetID: rand.Intn(numAssets),
			Body:    RandStringRunes(20),
		}
		if rand.Intn(2) == 1 {
			comment.ParentID = rand.Intn(currentNum)
		}
		payload, err := json.Marshal(comment)
		if err != nil {
			errors.Wrap(err, "Could not marshal JSON comment")
		}
		job := Job{
			Data: payload,
			Type: "comment",
		}
		jobs <- job
		currentNum++
	}
	return nil
}

// RandStringRunes generates a random string for our dummy data
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

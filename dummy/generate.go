package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
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
		currentNum++
		if err := sendToSponge(user, "user"); err != nil {
			return errors.Wrap(err, "Could not send item to sponge")
		}
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
		currentNum++
		if err := sendToSponge(asset, "asset"); err != nil {
			return errors.Wrap(err, "Could not send item to sponge")
		}
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
		currentNum++
		if err := sendToSponge(comment, "comment"); err != nil {
			return errors.Wrap(err, "Could not send item to sponge")
		}
	}
	return nil
}

// sendToSponge sends a generated item to sponge for processing.
func sendToSponge(item interface{}, typeIn string) error {

	url := fmt.Sprintf("http://%s/1.0/item/coral_%s", spongeHost, typeIn)
	payload, err := json.Marshal(item)
	if err != nil {
		return errors.Wrap(err, "Could not encode JSON payload to sponge")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrap(err, "Could not create http request")
	}
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "Could not execute POST request to sponge")
	}
	if res.StatusCode != http.StatusOK {
		return errors.Wrap(err, "Unexpected sponge response")
	}
	defer res.Body.Close()

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

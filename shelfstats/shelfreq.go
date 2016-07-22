package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Results includes all the results for each view provided by the shelf POC server.
type Results struct {
	UserComments       Comparison
	UserAssets         Comparison
	AssetUsersComments Comparison
	Parents            Comparison
	GrandParents       Comparison
	GrGrandParents     Comparison
}

// Comparison includes average time for Mongo and Graph requests respectively.
type Comparison struct {
	Mongo float64
	Graph float64
}

// Item is used to unmarshal item responses.
type Item struct {
	ID string `json:"id"`
}

// getAsset gets a random asset from the shelf API.
func getAsset() (string, error) {

	// Form the URL.
	numAssets := numDocs / 20
	url := fmt.Sprintf("http://%s:8080/asset?num=%d", shelfHost, numAssets)

	// Generate the request.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "Could not generate request to GET asset")
	}
	req.Header.Add("content-type", "application/json")

	// Send the request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Could not execute GET asset request")
	}
	defer res.Body.Close()

	// Make sure the retrieval was successful.
	if res.StatusCode != http.StatusOK {
		return "", errors.Wrap(err, "Received an unexpected response from shelf")
	}

	return getID(res)
}

// getComment gets a random comment from the shelf API.
func getComment() (string, error) {

	// Form the URL.
	numComments := (numDocs * 4) / 5
	url := fmt.Sprintf("http://%s:8080/comment?num=%d", shelfHost, numComments)

	// Generate the request.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "Could not generate request to GET comment")
	}
	req.Header.Add("content-type", "application/json")

	// Send the request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Could not execute GET comment request")
	}
	defer res.Body.Close()

	// Make sure the retrieval was successful.
	if res.StatusCode != http.StatusOK {
		return "", errors.Wrap(err, "Received an unexpected response from shelf")
	}

	return getID(res)
}

// getUser gets a random user from the shelf API.
func getUser() (string, error) {

	// Form the URL.
	numUsers := (numDocs * 3) / 20
	url := fmt.Sprintf("http://%s:8080/user?num=%d", shelfHost, numUsers)

	// Generate the request.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "Could not generate request to GET user")
	}
	req.Header.Add("content-type", "application/json")

	// Send the request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Could not execute GET user request")
	}
	defer res.Body.Close()

	// Make sure the retrieval was successful.
	if res.StatusCode != http.StatusOK {
		return "", errors.Wrap(err, "Received an unexpected response from shelf")
	}

	return getID(res)
}

// getID extracts the ID from an item response.
func getID(res *http.Response) (string, error) {
	var item Item
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Could not read response body")
	}
	if err = json.Unmarshal(body, &item); err != nil {
		return "", errors.Wrap(err, "Could not unmarshal JSON in response")
	}
	return item.ID, nil
}

// validate ensures that we have returned at least one item.
func validate(res *http.Response) (bool, error) {
	var items []Item
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, errors.Wrap(err, "Could not read response body")
	}
	if err = json.Unmarshal(body, &items); err != nil {
		return false, errors.Wrap(err, "Could not unmarshal JSON in response")
	}
	if len(items) == 0 {
		return false, nil
	}
	return true, nil
}

// timeResponse times an HTTP response from shelf in seconds (e.g., retrieving a view).
func timeResponse(url string, val bool) (float64, error) {

	// Generate the request.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0.0, errors.Wrap(err, "Could not generate a GET request")
	}
	req.Header.Add("content-type", "application/json")

	// Time the response.
	start := time.Now()
	res, err := http.DefaultClient.Do(req)
	elasped := time.Since(start)
	if err != nil {
		return elasped.Seconds(), errors.Wrap(err, "Could not execute GET request")
	}
	defer res.Body.Close()

	// Make sure the retrieval was successful.
	if res.StatusCode != http.StatusOK {
		return elasped.Seconds(), errors.Wrap(err, "Received an unexpected response from shelf")
	}

	// Make sure we actually got results, if validation is requested.
	if val {
		if results, err := validate(res); !results || err != nil {
			return -1.0, fmt.Errorf("No results")
		}
	}

	return elasped.Seconds(), nil
}

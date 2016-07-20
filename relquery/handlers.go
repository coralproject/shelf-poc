package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// GetAsset returns a random asset.
func GetAsset(w http.ResponseWriter, r *http.Request) {

	// Get the number of asset provided inthe query string.
	queryvals := r.URL.Query()
	numString := queryvals["num"][0]
	num, err := strconv.Atoi(numString)
	if err != nil {
		err = errors.Wrap(err, "Could not parse number of assets")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query mongo for a random asset.
	asset, err := retrieveRandAsset(num)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve random asset")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(asset); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// GraphQuery returns comment and author items corresponding to an asset
// using graphed relationships managed via Cayley.
func GraphQuery(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	assetID := queryvals["asset"][0]

	// Query cayley to get the item IDs related to this asset ID.
	itemIDs, err := getItemsOnAsset(assetID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve item IDs from Cayley.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query MongoDB to retrieve the corresponding documents.
	items, err := retrieveObjectList(itemIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve items from Mongo.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// MongoQuery returns comment and author items corresponding to an asset
// using embedded relationships in MongoDB.
func MongoQuery(w http.ResponseWriter, r *http.Request) {

	//queryvals := r.URL.Query()
	//assetID := queryvals["userid"][0]

	fmt.Fprint(w, "Welcome to this example API\n")
}

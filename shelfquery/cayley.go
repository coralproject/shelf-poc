package main

import (
	"github.com/cayleygraph/cayley"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/pkg/errors"
)

// openCayley opens a new connection to the Cayley graph DB.
func openCayley() (*cayley.Handle, error) {
	store, err := cayley.NewGraph("mongo", "localhost:27017", nil)
	if err != nil {
		return store, errors.Wrap(err, "Could not open and use the Cayley DB")
	}
	return store, nil
}

// getItemsOnAsset gets all the comments and authors related to an asset.
func getItemsOnAsset(assetID string) ([]string, error) {

	// Connect to cayley.
	store, err := openCayley()
	if err != nil {
		err = errors.Wrap(err, "Could not open connection to Cayley")
		return nil, err
	}
	defer store.Close()

	// Get the related item IDs.
	it, _ := cayley.StartPath(store, assetID).In("contextualized_with").Or(cayley.StartPath(store, assetID).In("contextualized_with").Out("authored_by")).BuildIterator().Optimize()
	defer it.Close()

	// Gather the results.
	var ids []string
	for cayley.RawNext(it) {
		if it.Result() != nil {
			ids = append(ids, store.NameOf(it.Result()))
		}
	}
	if it.Err() != nil {
		return nil, errors.Wrap(it.Err(), "Lost connection to Cayley")
	}

	return ids, nil
}

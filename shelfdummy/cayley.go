package main

import (
	"log"
	"time"

	"github.com/cayleygraph/cayley"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/pkg/errors"
)

func openCayley() (*cayley.Handle, error) {
	store, err := cayley.NewGraph("mongo", "localhost:27017", nil)
	if err != nil {
		return store, errors.Wrap(err, "Could not open and use the Cayley DB")
	}
	return store, nil
}

func applyTx() error {

	// Get the connection to Cayley.
	store, err := openCayley()
	if err != nil {
		return errors.Wrap(err, "Could not open connection to Cayley")
	}

	// Apply all the transactions that we have accumulated.
	txMutex.Lock()
	if err := store.ApplyTransaction(tx); err != nil {
		return errors.Wrap(err, "Could not execute transaction")
	}
	txMutex.Unlock()
	return nil
}

// applyPerTx applys transactions when they reach 100,000 quads.
func applyPerTx() {
	go func() {
		for {
			txMutex.Lock()
			if txCount >= 100000 {
				log.Println("Applying Cayley bulk transaction")
				err := applyTx()
				results <- err
				tx = cayley.NewTransaction()
				txCount = 0
			}
			txMutex.Unlock()
			time.Sleep(3 * time.Minute)
		}
	}()
}

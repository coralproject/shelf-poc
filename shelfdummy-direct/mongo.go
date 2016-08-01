package main

import (
	"math/rand"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Set of error variables.
var (
	ErrNotFound = errors.New("Set Not found")
)

// mongoCollections opens a new Mongo master session and return a session copy.
func mongoConnection(dbName string) (*db.DB, error) {

	// Define the mongo config.
	var mgoDB *db.DB
	config := mongo.Config{
		Host: mongoHost,
		DB:   dbName,
	}

	// Register the Mongo master session.
	err := db.RegMasterSession("Mongo", "test", config)
	if err != nil {
		return mgoDB, errors.Wrap(err, "Could not register Master Session")
	}

	// Create a session copy.
	mgoDB, err = db.NewMGO("Mongo", "test")
	if err != nil {
		return mgoDB, errors.Wrap(err, "Could not connect to MongoDB")
	}

	return mgoDB, err
}

// retrieveRandAsset retrieves a random item of type "coral_asset" from Mongo.
func retrieveRandAsset(num int, mgoDB *db.DB) (Item, error) {

	// Define the query.
	var result Item
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"t": "coral_asset"}).Limit(1).Skip(rand.Intn(num)).One(&result)
	}

	// Execute the query.
	if err := mgoDB.ExecuteMGO("Mongo", "coral_items", f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return Item{}, errors.Wrap(err, "Could not locate asset")
	}

	return result, nil
}

// retrieveRandUser retrieves a random item of type "coral_user" from Mongo.
func retrieveRandUser(num int, mgoDB *db.DB) (Item, error) {

	// Define the query.
	var result Item
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"t": "coral_user"}).Limit(1).Skip(rand.Intn(num)).One(&result)
	}

	// Execute the query.
	if err := mgoDB.ExecuteMGO("Mongo", "coral_items", f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return Item{}, errors.Wrap(err, "Could not locate user")
	}

	return result, nil
}

// retrieveRandComment retrieves a random item of type "coral_comment" from Mongo.
func retrieveRandComment(num int, mgoDB *db.DB) (Item, error) {

	// Define the query.
	var result Item
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"t": "coral_comment"}).Limit(1).Skip(rand.Intn(num)).One(&result)
	}

	// Execute the query.
	if err := mgoDB.ExecuteMGO("Mongo", "coral_items", f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return Item{}, errors.Wrap(err, "Could not locate comment")
	}

	return result, nil
}

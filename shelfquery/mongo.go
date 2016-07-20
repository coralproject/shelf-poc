package main

import (
	"math/rand"
	"time"

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

func init() {
	rand.Seed(time.Now().UnixNano())
}

// mongoCollections opens a new Mongo master session and return a session copy.
func mongoConnection(dbName string) (*db.DB, error) {

	// Define the mongo config.
	var mgoDB *db.DB
	config := mongo.Config{
		Host: "localhost",
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
func retrieveRandAsset(num int) (Item, error) {

	// Create a session copy.
	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return Item{}, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

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

// retrieveObjectList retrieves the items corresponding to the input list of object IDs.
func retrieveObjectList(ids []string) ([]Item, error) {

	// Create a session copy.
	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

	// Convert the string object IDs to bson.ObjectId.
	var idsBSON []bson.ObjectId
	for _, idString := range ids {
		idsBSON = append(idsBSON, bson.ObjectIdHex(idString))
	}

	// Form the query.
	var results []Item
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": idsBSON}}).All(&results)
	}

	// Execute the query.
	if err := mgoDB.ExecuteMGO("Mongo", "coral_items", f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return nil, errors.Wrap(err, "Could not locate items")
	}

	return results, nil
}

// retrieveCommentsByAsset retrieves comments "contextualized on" an asset.
func retrieveCommentsByAsset(assetID string) ([]Item, error) {

	// Create a session copy.
	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

	// Form the query.
	var results []Item
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"t": "coral_comment", "rels": bson.M{"$elemMatch": bson.M{"id": assetID}}}).All(&results)
	}

	// Execute the query.
	if err := mgoDB.ExecuteMGO("Mongo", "coral_items", f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return nil, errors.Wrap(err, "Could not locate items")
	}

	return results, nil
}

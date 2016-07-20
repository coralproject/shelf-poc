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

func mongoConnection(dbName string) (*db.DB, error) {

	var mgoDB *db.DB
	config := mongo.Config{
		Host: "localhost",
		DB:   dbName,
	}

	err := db.RegMasterSession("Mongo", "test", config)
	if err != nil {
		return mgoDB, errors.Wrap(err, "Could not register Master Session")
	}

	mgoDB, err = db.NewMGO("Mongo", "test")
	if err != nil {
		return mgoDB, errors.Wrap(err, "Could not connect to MongoDB")
	}

	return mgoDB, err

}

func retrieveRandAsset(num int) (Item, error) {

	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return Item{}, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

	var result Item
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"t": "coral_asset"}).Limit(1).Skip(rand.Intn(num)).One(&result)
	}

	if err := mgoDB.ExecuteMGO("Mongo", "coral_items", f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return Item{}, errors.Wrap(err, "Could not locate asset")
	}

	return result, nil
}

func retrieveObjectList(ids []string) ([]Item, error) {

	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

	var idsBSON []bson.ObjectId
	for _, idString := range ids {
		idsBSON = append(idsBSON, bson.ObjectIdHex(idString))
	}

	var results []Item
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"_id": bson.M{"$in": idsBSON}}).All(&results)
	}

	if err := mgoDB.ExecuteMGO("Mongo", "coral_items", f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return nil, errors.Wrap(err, "Could not locate items")
	}

	return results, nil
}

func retrieveCommentsByAsset(assetID string) ([]Item, error) {

	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

	var results []Item
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"t": "coral_comment", "rels": bson.M{"$elemMatch": bson.M{"id": assetID}}}).All(&results)
	}

	if err := mgoDB.ExecuteMGO("Mongo", "coral_items", f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return nil, errors.Wrap(err, "Could not locate items")
	}

	return results, nil
}

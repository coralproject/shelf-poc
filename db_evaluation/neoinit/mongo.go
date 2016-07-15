package main

import (
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

func mongoConnection(dbName string) (*db.DB, error) {

	var mgoDB *db.DB
	config := mongo.Config{
		Host: "localhost",
		DB:   "coral-poc",
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

type mongoDoc struct {
	UnderScoreID bson.ObjectId `bson:"_id,omitempty"`
	Type         string        `json:"type"`
}

type mongoDocs []mongoDoc

func getObjectIDs(db *db.DB, context interface{}, collection string) (mongoDocs, error) {

	var results mongoDocs
	f := func(c *mgo.Collection) error {
		return c.Find(nil).All(&results)
	}

	if err := db.ExecuteMGO(context, collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return nil, err
	}

	return results, nil

}

type Author struct {
	UnderScoreId bson.ObjectId `bson:"_id,omitempty"`
	Name         string        `json:"name"`
}

type Authors []Author

type Comment struct {
	UnderScoreId bson.ObjectId `bson:"_id,omitempty"`
	Author       string        `json:"author"`
}

type Comments []Comment

func getAuthors(mgoDB *db.DB, context interface{}, collection string) (Authors, error) {

	var results Authors
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"type": "author"}).All(&results)
	}

	if err := mgoDB.ExecuteMGO(context, collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return nil, err
	}

	return results, nil

}

func getComments(mgoDB *db.DB, context interface{}, collection string) (Comments, error) {

	var results Comments
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"type": "comment"}).All(&results)
	}

	if err := mgoDB.ExecuteMGO(context, collection, f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return nil, err
	}

	return results, nil

}

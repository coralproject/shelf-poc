package main

import (
	"fmt"
	"strings"

	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
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

// saveView retrieves the items corresponding to the input list of object IDs
// and saves the view to another Mongo collection.
func saveView(ids []string, fillRels []FillRel) (CollectionOut, error) {

	// Create a session copy.
	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return CollectionOut{}, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

	// Convert the string object IDs to bson.ObjectId.
	var idsBSON []bson.ObjectId
	for _, idString := range ids {
		idsBSON = append(idsBSON, bson.ObjectIdHex(idString))
	}

	// Generate the collection name.
	uid := uuid.NewV4()
	newCol := fmt.Sprintf("%s%s", "view", uid)
	newCol = strings.Replace(newCol, "-", "", -1)

	// Form the query.
	var results []Item
	f := func(c *mgo.Collection) error {
		return c.Pipe([]bson.M{{"$match": bson.M{"_id": bson.M{"$in": idsBSON}}}, {"$out": newCol}}).All(&results)
	}

	// Execute the query.
	if err := mgoDB.ExecuteMGO("Mongo", "coral_items", f); err != nil {
		if err == mgo.ErrNotFound {
			err = ErrNotFound
		}
		return CollectionOut{}, errors.Wrap(err, "Could not locate items")
	}

	// Form the response.
	out := CollectionOut{
		Number:     len(ids),
		Collection: newCol,
	}

	// Queue the relationships to be embedded.
	for _, rel := range fillRels {
		rel.Collection = newCol
		fills <- rel
	}

	return out, nil
}

// fillRelationships embeds relationships into saved views.
func fillRelationships() {
	go func() {
		for rel := range fills {
			// Create a session copy.
			mgoDB, err := db.NewMGO("Mongo", "test")
			if err != nil {
				panic(errors.Wrap(err, "Could not connect to MongoDB"))
			}

			// Insert or update the item.
			f := func(c *mgo.Collection) error {
				q := bson.M{"_id": rel.ID}
				update := bson.M{"$push": bson.M{"rels": rel.Relationship}}
				err := c.Update(q, update)
				return err
			}

			if err := mgoDB.ExecuteMGO("Mongo", rel.Collection, f); err != nil {
				panic(err)
			}
			mgoDB.CloseMGO("Mongo")
		}
	}()
}

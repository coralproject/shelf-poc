package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

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

// retrieveRandUser retrieves a random item of type "coral_user" from Mongo.
func retrieveRandUser(num int) (Item, error) {

	// Create a session copy.
	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return Item{}, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

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
func retrieveRandComment(num int) (Item, error) {

	// Create a session copy.
	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return Item{}, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

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

// retrieveCommentsByUser retrieves comments authored by the given user.
func retrieveCommentsByUser(userID string) ([]Item, error) {

	// Create a session copy.
	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

	// Form the query.
	var results []Item
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"t": "coral_comment", "rels": bson.M{"$elemMatch": bson.M{"id": userID}}}).All(&results)
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

// retrieveCommentsByParents retrieves comments parented by one of the given list of comments.
func retrieveCommentsByParents(parentIDs []string) ([]Item, error) {

	// Create a session copy.
	mgoDB, err := db.NewMGO("Mongo", "test")
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to MongoDB")
	}
	defer mgoDB.CloseMGO("Mongo")

	// Form the query.
	var results []Item
	f := func(c *mgo.Collection) error {
		return c.Find(bson.M{"t": "coral_comment", "rels": bson.M{"$elemMatch": bson.M{"id": bson.M{"$in": parentIDs}}}}).All(&results)
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

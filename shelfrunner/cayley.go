package main

import (
	"fmt"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/cayleygraph/cayley/quad"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

// openCayley opens a new connection to the Cayley graph DB.
func openCayley() (*cayley.Handle, error) {
	host := fmt.Sprintf("%s:27017", mongoHost)
	store, err := cayley.NewGraph("mongo", host, nil)
	if err != nil {
		return store, errors.Wrap(err, "Could not open and use the Cayley DB")
	}
	return store, nil
}

// getItemsOnAsset gets all the comments and authors related to an asset.
func getItemsOnAsset(assetID string) ([]string, []FillRel, error) {

	// Connect to cayley.
	store, err := openCayley()
	if err != nil {
		err = errors.Wrap(err, "Could not open connection to Cayley")
		return nil, nil, err
	}
	defer store.Close()

	// Get the related item IDs and relationships.
	path := cayley.StartPath(store, quad.String(assetID)).In(quad.String("contextualized_with")).Tag("comment").In(quad.String("authored")).Tag("user")

	it := path.BuildIterator()
	it, _ = it.Optimize()
	defer it.Close()

	var ids []string
	var rels []FillRel
	for it.Next() {
		tags := make(map[string]graph.Value)
		it.TagResults(tags)
		if t1, ok := tags["comment"]; ok {
			commentID := quad.NativeOf(store.NameOf(t1)).(string)
			ids = append(ids, commentID)
			rel := FillRel{
				ID: bson.ObjectIdHex(commentID),
				Relationship: Rel{
					Name: "context",
					Type: "coral_asset",
					ID:   assetID,
				},
			}
			rels = append(rels, rel)
			if t2, ok := tags["user"]; ok {
				rel := FillRel{
					ID: bson.ObjectIdHex(commentID),
					Relationship: Rel{
						Name: "author",
						Type: "coral_user",
						ID:   quad.NativeOf(store.NameOf(t2)).(string),
					},
				}
				rels = append(rels, rel)
			}
		}
		if it.Result() != nil {
			ids = append(ids, quad.NativeOf(store.NameOf(it.Result())).(string))
		}
	}

	if it.Err() != nil {
		return nil, nil, errors.Wrap(it.Err(), "Lost connection to Cayley")
	}

	return ids, rels, nil
}

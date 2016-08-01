package main

import (
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/cayleygraph/cayley/quad"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
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

// getAssetsOnUser gets all the assets related to a user.
func getAssetsOnUser(userID string) ([]string, []FillRel, error) {

	// Connect to cayley.
	store, err := openCayley()
	if err != nil {
		err = errors.Wrap(err, "Could not open connection to Cayley")
		return nil, nil, err
	}
	defer store.Close()

	// Get the related item IDs.
	path := cayley.StartPath(store, quad.String(userID)).Out(quad.String("authored")).Tag("comment").Out(quad.String("contextualized_with")).Tag("asset")

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
			rel1 := FillRel{
				Relationship: Rel{
					Name: "comment",
					Type: "coral_comment",
					ID:   commentID,
				},
			}
			if t2, ok := tags["asset"]; ok {
				rel1.ID = bson.ObjectIdHex(quad.NativeOf(store.NameOf(t2)).(string))
				rel2 := FillRel{
					ID: bson.ObjectIdHex(quad.NativeOf(store.NameOf(t2)).(string)),
					Relationship: Rel{
						Name: "author",
						Type: "coral_user",
						ID:   userID,
					},
				}
				rels = append(rels, rel1)
				rels = append(rels, rel2)
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

// getCommentsOnUser gets all the comments related to a user.
func getCommentsOnUser(userID string) ([]string, []FillRel, error) {

	// Connect to cayley.
	store, err := openCayley()
	if err != nil {
		err = errors.Wrap(err, "Could not open connection to Cayley")
		return nil, nil, err
	}
	defer store.Close()

	// Get the related item IDs.
	path := cayley.StartPath(store, quad.String(userID)).Out(quad.String("authored")).Tag("comment")

	it := path.BuildIterator()
	it, _ = it.Optimize()
	defer it.Close()

	var ids []string
	var rels []FillRel
	for it.Next() {
		tags := make(map[string]graph.Value)
		it.TagResults(tags)
		if t, ok := tags["comment"]; ok {
			commentID := quad.NativeOf(store.NameOf(t)).(string)
			rel := FillRel{
				ID: bson.ObjectIdHex(commentID),
				Relationship: Rel{
					Name: "author",
					Type: "coral_user",
					ID:   userID,
				},
			}
			rels = append(rels, rel)
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

// getCommentsOnPar gets all the comments parented by comments that are authored by
// the author of the parent of the comment provided.
func getCommentsOnPar(commentID string) ([]string, []FillRel, error) {

	// Connect to cayley.
	store, err := openCayley()
	if err != nil {
		err = errors.Wrap(err, "Could not open connection to Cayley")
		return nil, nil, err
	}
	defer store.Close()

	// Get the related item IDs.
	path := cayley.StartPath(store, quad.String(commentID)).Out(quad.String("parented_by")).In(quad.String("authored")).Out(quad.String("authored")).Tag("parent").In(quad.String("parented_by")).Tag("child")

	it := path.BuildIterator()
	it, _ = it.Optimize()
	defer it.Close()

	var ids []string
	var rels []FillRel
	for it.Next() {
		tags := make(map[string]graph.Value)
		it.TagResults(tags)
		if t1, ok := tags["parent"]; ok {
			parentID := quad.NativeOf(store.NameOf(t1)).(string)
			if t2, ok := tags["child"]; ok {
				rel := FillRel{
					ID: bson.ObjectIdHex(quad.NativeOf(store.NameOf(t2)).(string)),
					Relationship: Rel{
						Name: "parent",
						Type: "coral_comment",
						ID:   parentID,
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

// getGrandCommentsOnPar gets all the comments grandparented by comments that are authored by
// the author of the parent of the comment provided.
func getGrandCommentsOnPar(commentID string) ([]string, []FillRel, error) {

	// Connect to cayley.
	store, err := openCayley()
	if err != nil {
		err = errors.Wrap(err, "Could not open connection to Cayley")
		return nil, nil, err
	}
	defer store.Close()

	// Get the related item IDs.
	path := cayley.StartPath(store, quad.String(commentID)).Out(quad.String("parented_by")).In(quad.String("authored")).Out(quad.String("authored")).Tag("parent").In(quad.String("parented_by")).Tag("child").In(quad.String("parented_by")).Tag("grandchild")

	it := path.BuildIterator()
	it, _ = it.Optimize()
	defer it.Close()

	var ids []string
	var rels []FillRel
	for it.Next() {
		tags := make(map[string]graph.Value)
		it.TagResults(tags)
		if t1, ok := tags["parent"]; ok {
			grandParentID := quad.NativeOf(store.NameOf(t1)).(string)
			if t2, ok := tags["child"]; ok {
				parentID := quad.NativeOf(store.NameOf(t2)).(string)
				if t3, ok := tags["grandchild"]; ok {
					rel := FillRel{
						ID: bson.ObjectIdHex(quad.NativeOf(store.NameOf(t3)).(string)),
						Relationship: Rel{
							Name: "parent",
							Type: "coral_comment",
							ID:   parentID,
						},
					}
					rels = append(rels, rel)
					rel = FillRel{
						ID: bson.ObjectIdHex(quad.NativeOf(store.NameOf(t3)).(string)),
						Relationship: Rel{
							Name: "grandparent",
							Type: "coral_comment",
							ID:   grandParentID,
						},
					}
					rels = append(rels, rel)
				}
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

// getGrGrandCommentsOnPar gets all the comments great-grandparented by comments that
// are authored by the author of the parent of the comment provided.
func getGrGrandCommentsOnPar(commentID string) ([]string, []FillRel, error) {

	// Connect to cayley.
	store, err := openCayley()
	if err != nil {
		err = errors.Wrap(err, "Could not open connection to Cayley")
		return nil, nil, err
	}
	defer store.Close()

	// Get the related item IDs.
	path := cayley.StartPath(store, quad.String(commentID)).Out(quad.String("parented_by")).In(quad.String("authored")).Out(quad.String("authored")).Tag("parent").In(quad.String("parented_by")).Tag("child").In(quad.String("parented_by")).Tag("grandchild").In(quad.String("parented_by")).Tag("greatgrandchild")

	it := path.BuildIterator()
	it, _ = it.Optimize()
	defer it.Close()

	var ids []string
	var rels []FillRel
	for it.Next() {
		tags := make(map[string]graph.Value)
		it.TagResults(tags)
		if t1, ok := tags["parent"]; ok {
			greatGrandParentID := quad.NativeOf(store.NameOf(t1)).(string)
			if t2, ok := tags["child"]; ok {
				grandparentID := quad.NativeOf(store.NameOf(t2)).(string)
				if t3, ok := tags["grandchild"]; ok {
					parentID := quad.NativeOf(store.NameOf(t3)).(string)
					if t4, ok := tags["greatgrandchild"]; ok {
						rel := FillRel{
							ID: bson.ObjectIdHex(quad.NativeOf(store.NameOf(t4)).(string)),
							Relationship: Rel{
								Name: "parent",
								Type: "coral_comment",
								ID:   parentID,
							},
						}
						rels = append(rels, rel)
						rel = FillRel{
							ID: bson.ObjectIdHex(quad.NativeOf(store.NameOf(t4)).(string)),
							Relationship: Rel{
								Name: "grandparent",
								Type: "coral_comment",
								ID:   grandparentID,
							},
						}
						rels = append(rels, rel)
						rel = FillRel{
							ID: bson.ObjectIdHex(quad.NativeOf(store.NameOf(t4)).(string)),
							Relationship: Rel{
								Name: "greatgrandparent",
								Type: "coral_comment",
								ID:   greatGrandParentID,
							},
						}
						rels = append(rels, rel)

					}
				}
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

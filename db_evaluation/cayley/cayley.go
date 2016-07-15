package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/mongo"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

func openCayley() (*cayley.Handle, error) {
	// Initialize the database
	if err := graph.InitQuadStore("mongo", "localhost:27017", nil); err != nil {
		return nil, errors.Wrap(err, "Could not initialize quad store")
	}

	// Open and use the database
	store, err := cayley.NewGraph("mongo", "localhost:27017", nil)
	if err != nil {
		return store, errors.Wrap(err, "Could not open and use the Cayley DB")
	}
	return store, nil
}

func storeNodes(store *cayley.Handle, items mongoDocs) error {

	t := cayley.NewTransaction()
	for _, item := range items {
		t.AddQuad(cayley.Quad(item.UnderScoreID.Hex(), "is_type", item.Type, ""))
	}

	if err := store.ApplyTransaction(t); err != nil {
		return errors.Wrap(err, "Could not execute transaction")
	}
	return nil
}

func authorRelationships(store *cayley.Handle, mgoDB *db.DB, con interface{}, col string) error {

	authors, err := getAuthors(mgoDB, con, col)
	if err != nil {
		return errors.Wrap(err, "Could not get authors from MongoDB")
	}

	comments, err := getComments(mgoDB, con, col)
	if err != nil {
		return errors.Wrap(err, "Could not get comments from MongoDB")
	}

	t := cayley.NewTransaction()
	for _, comment := range comments {
		authorName := comment.Author
		for _, author := range authors {
			if author.Name == authorName {
				t.AddQuad(cayley.Quad(comment.UnderScoreId.Hex(), "authored_by",
					author.UnderScoreId.Hex(), ""))
			}
		}
	}
	if err := store.ApplyTransaction(t); err != nil {
		return errors.Wrap(err, "Could not apply authored_by transactions")
	}

	return nil

}

func createThread(store *cayley.Handle, name string) (string, error) {

	threadID := uuid.NewV4()
	threadIDString := fmt.Sprintf("%s", threadID)

	t := cayley.NewTransaction()
	t.AddQuad(cayley.Quad(threadIDString, "is_type", "thread", ""))
	t.AddQuad(cayley.Quad(threadIDString, "is_named", name, ""))

	if err := store.ApplyTransaction(t); err != nil {
		return threadIDString, errors.Wrap(err, "Could not apply thread transactions")
	}
	return threadIDString, nil
}

func addSomeToThread(store *cayley.Handle, docs mongoDocs, threadID string) error {

	t := cayley.NewTransaction()
	for _, doc := range docs {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		if r1.Intn(2) == 1 {
			t.AddQuad(cayley.Quad(doc.UnderScoreID.Hex(), "threaded_on", threadID, ""))
		}
	}

	if err := store.ApplyTransaction(t); err != nil {
		return errors.Wrap(err, "Could not apply threaded_on transactions")
	}
	return nil
}

func getItemsOnThread(store *cayley.Handle, threadID string) error {
	it, _ := cayley.StartPath(store, threadID).In("threaded_on").BuildIterator().Optimize()
	defer it.Close()
	for cayley.RawNext(it) {
		if it.Result() != nil {
			fmt.Println(store.NameOf(it.Result()))
		}
	}
	if it.Err() != nil {
		return errors.Wrap(it.Err(), "Lost connection to Cayley")
	}
	return nil
}

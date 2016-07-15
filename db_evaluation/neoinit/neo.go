package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/ardanlabs/kit/db"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	_ "gopkg.in/cq.v1"
)

func neoConnection() (*sql.DB, error) {
	return sql.Open("neo4j-cypher", "http://localhost:7474")
}

func insertNodes(items mongoDocs, neoDB *sql.DB) error {

	for _, item := range items {

		stmt, err := neoDB.Prepare("create (:Item {id: {0}, type: {1}})")
		if err != nil {
			return errors.Wrap(err, "Could not prepare create statement")
		}

		_, err = stmt.Exec(item.UnderScoreID.Hex(), item.Type)
		stmt.Close()
		if err != nil {
			return errors.Wrap(err, "Could not execute prepared statement")
		}

	}

	return nil

}

func authorRelationships(neoDB *sql.DB, mgoDB *db.DB, con interface{}, col string) error {

	authors, err := getAuthors(mgoDB, con, col)
	if err != nil {
		return errors.Wrap(err, "Could not get authors from MongoDB")
	}

	comments, err := getComments(mgoDB, con, col)
	if err != nil {
		return errors.Wrap(err, "Could not get comments from MongoDB")
	}

	for _, comment := range comments {

		authorName := comment.Author
		for _, author := range authors {
			if author.Name == authorName {

				stmt, err := neoDB.Prepare(`
				MATCH (i1:Item {id:{0}}), (i2:Item {id:{1}})
				CREATE (i1)-[:AUTHORED_BY]->(i2)
				`)
				if err != nil {
					return errors.Wrap(err, "Could not prepare create stmt")
				}

				_, err = stmt.Exec(comment.UnderScoreId.Hex(), author.UnderScoreId.Hex())
				stmt.Close()
				if err != nil {
					return errors.Wrap(err, "Could not execute stmt")
				}

			}
		}

	}

	return nil

}

func createThread(neoDB *sql.DB, name string) (string, error) {

	threadID := uuid.NewV4()
	threadIDString := fmt.Sprintf("%s", threadID)

	stmt, err := neoDB.Prepare("create (:Thread {id: {0}, name: {1}})")
	if err != nil {
		return threadIDString, errors.Wrap(err, "Could not prepare create statement")
	}

	_, err = stmt.Exec(threadIDString, name)
	stmt.Close()
	if err != nil {
		return threadIDString, errors.Wrap(err, "Could not execute prepared statement")
	}

	return threadIDString, nil

}

func addSomeToThread(neoDB *sql.DB, docs mongoDocs, threadID string) error {

	for _, doc := range docs {

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		if r1.Intn(2) == 1 {

			stmt, err := neoDB.Prepare(`
				MATCH (t:Thread {id:{0}}), (i:Item {id:{1}})
				CREATE (i)-[:THREADED_ON]->(t)
				`)
			if err != nil {
				return errors.Wrap(err, "Could not prepare create stmt")
			}

			_, err = stmt.Exec(threadID, doc.UnderScoreID.Hex())
			stmt.Close()
			if err != nil {
				return errors.Wrap(err, "Could not execute stmt")
			}

		}

	}

	return nil

}

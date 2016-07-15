package main

import (
	"log"
	"time"
)

func main() {

	time.Sleep(5 * time.Second)

	log.Println("Connecting to MongoDB")
	mgoDB, err := mongoConnection("coral-poc")
	if err != nil {
		log.Fatal(err)
	}
	defer mgoDB.CloseMGO("Mongo")

	log.Println("Get Object IDs from MongoDB")
	items, err := getObjectIDs(mgoDB, "Mongo", "Items")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connecting to Neo4j")
	neoDB, err := neoConnection()
	defer neoDB.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserting nodes into Neo4j")
	err = insertNodes(items, neoDB)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserting author relationships into Neo4j")
	err = authorRelationships(neoDB, mgoDB, "Mongo", "Items")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating a thread node")
	threadID, err := createThread(neoDB, "test_thread")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Adding items to the thread node")
	err = addSomeToThread(neoDB, items, threadID)
	if err != nil {
		log.Fatal(err)
	}

}

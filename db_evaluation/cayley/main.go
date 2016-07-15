package main

import "log"

func main() {

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

	log.Println("Create and connect to Cayley")
	store, err := openCayley()
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	log.Println("Store nodes in Cayley")
	if err = storeNodes(store, items); err != nil {
		log.Fatal(err)
	}

	log.Println("Inserting author relationships into Cayley")
	if err = authorRelationships(store, mgoDB, "Mongo", "Items"); err != nil {
		log.Fatal(err)
	}

	log.Println("Creating a thread node")
	threadID, err := createThread(store, "test_thread")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Adding items to the thread node")
	if err = addSomeToThread(store, items, threadID); err != nil {
		log.Fatal(err)
	}

	log.Println("Print out the items on the created test thread")
	if err = getItemsOnThread(store, threadID); err != nil {
		log.Fatal(err)
	}
}

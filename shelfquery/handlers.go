package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// GetAsset returns a random asset.
func GetAsset(w http.ResponseWriter, r *http.Request) {

	// Get the number of asset provided inthe query string.
	queryvals := r.URL.Query()
	numString := queryvals["num"][0]
	num, err := strconv.Atoi(numString)
	if err != nil {
		err = errors.Wrap(err, "Could not parse number of assets")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query mongo for a random asset.
	asset, err := retrieveRandAsset(num)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve random asset")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(asset); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// GetUser returns a random user.
func GetUser(w http.ResponseWriter, r *http.Request) {

	// Get the number of users provided inthe query string.
	queryvals := r.URL.Query()
	numString := queryvals["num"][0]
	num, err := strconv.Atoi(numString)
	if err != nil {
		err = errors.Wrap(err, "Could not parse number of users")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query mongo for a random user.
	user, err := retrieveRandUser(num)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve random user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// GetComment returns a random comment.
func GetComment(w http.ResponseWriter, r *http.Request) {

	// Get the number of comments provided in the query string.
	queryvals := r.URL.Query()
	numString := queryvals["num"][0]
	num, err := strconv.Atoi(numString)
	if err != nil {
		err = errors.Wrap(err, "Could not parse number of comments")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query mongo for a random comment.
	comment, err := retrieveRandComment(num)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve random parented comment")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// GraphQuerySingle returns comment and author items corresponding to an asset
// using graphed relationships managed via Cayley.
func GraphQuerySingle(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID and storage bool from the query string.
	queryvals := r.URL.Query()
	assetID := queryvals["asset"][0]
	saveIn := queryvals["save"][0]
	save, err := strconv.ParseBool(saveIn)
	if err != nil {
		err = errors.Wrap(err, "Could not parse save parameter.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query cayley to get the item IDs related to this asset ID.
	itemIDs, fillRels, err := getItemsOnAsset(assetID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve item IDs from Cayley.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If we are not persisting the view query MongoDB to retrieve
	// the corresponding documents out send those back in the response.
	if !save {
		items, err := retrieveObjectList(itemIDs)
		if err != nil {
			err = errors.Wrap(err, "Could not retrieve items from Mongo.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Encode the results.
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(items); err != nil {
			log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
		}
		return
	}

	// Otherwise, save the view to a Mongo collection.
	output, err := saveView(itemIDs, fillRels)
	if err != nil {
		err = errors.Wrap(err, "Could not save generated view.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return

}

// GraphQueryUserAssets returns all asset items commented on by the given user
// using graphed relationships managed via Cayley.
func GraphQueryUserAssets(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	userID := queryvals["user"][0]

	// Query cayley to get the item IDs related to this asset ID.
	itemIDs, err := getAssetsOnUser(userID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve item IDs from Cayley.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query MongoDB to retrieve the corresponding documents.
	items, err := retrieveObjectList(itemIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve items from Mongo.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// GraphQueryUserComments returns all comment items authored by the given user
// using graphed relationships managed via Cayley.
func GraphQueryUserComments(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	userID := queryvals["user"][0]

	// Query cayley to get the item IDs related to this asset ID.
	itemIDs, err := getCommentsOnUser(userID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve item IDs from Cayley.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query MongoDB to retrieve the corresponding documents.
	items, err := retrieveObjectList(itemIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve items from Mongo.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// GraphQueryParComments returns all comment items parented by comments authored by
// the author of the parent comment of the comment provided.
func GraphQueryParComments(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	commentID := queryvals["comment"][0]

	// Query cayley to get the item IDs related to this asset ID.
	itemIDs, err := getCommentsOnPar(commentID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve item IDs from Cayley.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query MongoDB to retrieve the corresponding documents.
	items, err := retrieveObjectList(itemIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve items from Mongo.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// GraphQueryGrandparComments returns all comment items grandparented by comments authored by
// the author of the parent comment of the comment provided.
func GraphQueryGrandparComments(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	commentID := queryvals["comment"][0]

	// Query cayley to get the item IDs related to this asset ID.
	itemIDs, err := getGrandCommentsOnPar(commentID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve item IDs from Cayley.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query MongoDB to retrieve the corresponding documents.
	items, err := retrieveObjectList(itemIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve items from Mongo.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// GraphQueryGrGrandparComments returns all comment items great-grandparented by
// comments authored by the author of the parent comment of the comment provided.
func GraphQueryGrGrandparComments(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	commentID := queryvals["comment"][0]

	// Query cayley to get the item IDs related to this asset ID.
	itemIDs, err := getGrGrandCommentsOnPar(commentID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve item IDs from Cayley.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query MongoDB to retrieve the corresponding documents.
	items, err := retrieveObjectList(itemIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve items from Mongo.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// MongoQueryUserAssets returns asset items commented on by the given user
// using embedded relationships in MongoDB.
func MongoQueryUserAssets(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	userID := queryvals["user"][0]

	// Get the comment items corresponding to the asset.
	comments, err := retrieveCommentsByUser(userID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comments from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the asset IDs from the comments.
	var assets []string
	for _, item := range comments {
		for _, rel := range item.Rels {
			if rel.Type == "coral_asset" {
				assets = append(assets, rel.ID)
			}
		}
	}

	// Get the author items corresponding to the extracted IDs.
	items, err := retrieveObjectList(assets)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve assets from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// MongoQueryUserComments returns comment items authored on by the given user
// using embedded relationships in MongoDB.
func MongoQueryUserComments(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	userID := queryvals["user"][0]

	// Get the comment items corresponding to the asset.
	comments, err := retrieveCommentsByUser(userID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comments from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// MongoQueryParComments returns comment items parented by comments authored
// by the author of the parent of the provided comment.
func MongoQueryParComments(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	commentID := queryvals["comment"][0]

	// Get the comment.
	comments, err := retrieveObjectList([]string{commentID})
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comment from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the parent comment ID.
	var parentID string
	for _, rel := range comments[0].Rels {
		if rel.Name == "parent" {
			parentID = rel.ID
		}
	}
	if parentID == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode([]Item{}); err != nil {
			log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
		}
		return
	}

	// Get the parent comment.
	comments, err = retrieveObjectList([]string{parentID})
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comment from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the author.
	var authorID string
	for _, rel := range comments[0].Rels {
		if rel.Name == "author" {
			authorID = rel.ID
		}
	}

	// Get all the comments authored by that author.
	comments, err = retrieveCommentsByUser(authorID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comment by user from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the IDs from the comments.
	var commentIDs []string
	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.ID.Hex())
	}

	// Get any child comments parented by these comment IDs.
	children, err := retrieveCommentsByParents(commentIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve children by parents from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(children); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// MongoQueryGrandparComments returns comment items grandparented by comments authored
// by the author of the parent of the provided comment.
func MongoQueryGrandparComments(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	commentID := queryvals["comment"][0]

	// Get the comment.
	comments, err := retrieveObjectList([]string{commentID})
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comment from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the parent comment ID.
	var parentID string
	for _, rel := range comments[0].Rels {
		if rel.Name == "parent" {
			parentID = rel.ID
		}
	}
	if parentID == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode([]Item{}); err != nil {
			log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
		}
		return
	}

	// Get the parent comment.
	comments, err = retrieveObjectList([]string{parentID})
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comment from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the author.
	var authorID string
	for _, rel := range comments[0].Rels {
		if rel.Name == "author" {
			authorID = rel.ID
		}
	}

	// Get all the comments authored by that author.
	comments, err = retrieveCommentsByUser(authorID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comment by user from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the IDs from the comments.
	var commentIDs []string
	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.ID.Hex())
	}

	// Get any child comments parented by these comment IDs.
	children, err := retrieveCommentsByParents(commentIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve children by parents from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(children) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode([]Item{}); err != nil {
			log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
		}
		return
	}
	commentIDs = []string{}
	for _, comment := range children {
		commentIDs = append(commentIDs, comment.ID.Hex())
	}

	// Get any child comments parented by these comment IDs.
	grandchildren, err := retrieveCommentsByParents(commentIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve children by parents from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(grandchildren); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// MongoQueryGrGrandparComments returns comment items great-grandparented by comments authored
// by the author of the parent of the provided comment.
func MongoQueryGrGrandparComments(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	commentID := queryvals["comment"][0]

	// Get the comment.
	comments, err := retrieveObjectList([]string{commentID})
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comment from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the parent comment ID.
	var parentID string
	for _, rel := range comments[0].Rels {
		if rel.Name == "parent" {
			parentID = rel.ID
		}
	}
	if parentID == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode([]Item{}); err != nil {
			log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
		}
		return
	}

	// Get the parent comment.
	comments, err = retrieveObjectList([]string{parentID})
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comment from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the author.
	var authorID string
	for _, rel := range comments[0].Rels {
		if rel.Name == "author" {
			authorID = rel.ID
		}
	}

	// Get all the comments authored by that author.
	comments, err = retrieveCommentsByUser(authorID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comment by user from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the IDs from the comments.
	var commentIDs []string
	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.ID.Hex())
	}

	// Get any child comments parented by these comment IDs.
	children, err := retrieveCommentsByParents(commentIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve children by parents from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(children) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode([]Item{}); err != nil {
			log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
		}
		return
	}
	commentIDs = []string{}
	for _, comment := range children {
		commentIDs = append(commentIDs, comment.ID.Hex())
	}

	// Get any child comments parented by these comment IDs.
	grandchildren, err := retrieveCommentsByParents(commentIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve children by parents from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(grandchildren) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode([]Item{}); err != nil {
			log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
		}
		return
	}
	commentIDs = []string{}
	for _, comment := range grandchildren {
		commentIDs = append(commentIDs, comment.ID.Hex())
	}

	// Get any child comments parented by these comment IDs.
	greatgrandchildren, err := retrieveCommentsByParents(commentIDs)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve children by parents from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(greatgrandchildren); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

// MongoQuerySingle returns comment and author items corresponding to an asset
// using embedded relationships in MongoDB.
func MongoQuerySingle(w http.ResponseWriter, r *http.Request) {

	// Get the asset ID from the query string.
	queryvals := r.URL.Query()
	assetID := queryvals["asset"][0]

	// Get the comment items corresponding to the asset.
	items, err := retrieveCommentsByAsset(assetID)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve comments from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the author IDs from the comments.
	var authors []string
	for _, item := range items {
		for _, rel := range item.Rels {
			if rel.Type == "coral_user" {
				authors = append(authors, rel.ID)
			}
		}
	}

	// Get the author items corresponding to the extracted IDs.
	authorItems, err := retrieveObjectList(authors)
	if err != nil {
		err = errors.Wrap(err, "Could not retrieve authors from MongoDB.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Concatenate the comments and authors.
	items = append(items, authorItems...)

	// Encode the results.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Printf("%s: %s", "ERROR Could not encode JSON response", err.Error())
	}
	return
}

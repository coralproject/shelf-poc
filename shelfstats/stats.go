package main

import (
	"fmt"
	"log"

	"github.com/gonum/stat"
	"github.com/pkg/errors"
)

const (
	userCommentsGraphURL    = "http://%s:8080/graph/usercomments?user=%s&save=false"
	userCommentsMongoURL    = "http://%s:8080/mongo/usercomments?user=%s"
	userAssetsGraphURL      = "http://%s:8080/graph/userassets?user=%s&save=false"
	userAssetsMongoURL      = "http://%s:8080/mongo/userassets?user=%s"
	assetGraphURL           = "http://%s:8080/graph/singleasset?asset=%s&save=false"
	assetMongoURL           = "http://%s:8080/mongo/singleasset?asset=%s"
	parentedGraphURL        = "http://%s:8080/graph/parentcomments?comment=%s&save=false"
	parentedMongoURL        = "http://%s:8080/mongo/parentcomments?comment=%s"
	grandParentedGraphURL   = "http://%s:8080/graph/grandparentcomments?comment=%s&save=false"
	grandParentedMongoURL   = "http://%s:8080/mongo/grandparentcomments?comment=%s"
	grGrandParentedGraphURL = "http://%s:8080/graph/greatgrandparentcomments?comment=%s&save=false"
	grGrandParentedMongoURL = "http://%s:8080/mongo/greatgrandparentcomments?comment=%s"
)

// userComments gathers stats on the user comments view requests.
func userComments() error {

	// Gather statistics.
	var mongoStats []float64
	var graphStats []float64
	for i := 1; i <= numRequests; i++ {
		user, err := getUser()
		if err != nil {
			return errors.Wrap(err, "Could not retrieve a random user")
		}
		graphurl := fmt.Sprintf(userCommentsGraphURL, shelfHost, user)
		mongourl := fmt.Sprintf(userCommentsMongoURL, shelfHost, user)
		graphStat, err := timeResponse(graphurl, false)
		if err != nil {
			return errors.Wrap(err, "Could not time graph request")
		}
		mongoStat, err := timeResponse(mongourl, false)
		if err != nil {
			return errors.Wrap(err, "Could not time mongorequest")
		}
		mongoStats = append(mongoStats, mongoStat)
		graphStats = append(graphStats, graphStat)
	}

	// Calculate means.
	mongoMean := stat.Mean(mongoStats, nil)
	graphMean := stat.Mean(graphStats, nil)
	comp := Comparison{
		Mongo: mongoMean,
		Graph: graphMean,
	}
	results.UserComments = comp
	return nil
}

// userAssets gathers stats on the user assets view requests.
func userAssets() error {

	// Gather statistics.
	var mongoStats []float64
	var graphStats []float64
	for i := 1; i <= numRequests; i++ {
		user, err := getUser()
		if err != nil {
			return errors.Wrap(err, "Could not retrieve a random user")
		}
		graphurl := fmt.Sprintf(userAssetsGraphURL, shelfHost, user)
		mongourl := fmt.Sprintf(userAssetsMongoURL, shelfHost, user)
		graphStat, err := timeResponse(graphurl, false)
		if err != nil {
			return errors.Wrap(err, "Could not time graph request")
		}
		mongoStat, err := timeResponse(mongourl, false)
		if err != nil {
			return errors.Wrap(err, "Could not time mongorequest")
		}
		mongoStats = append(mongoStats, mongoStat)
		graphStats = append(graphStats, graphStat)
	}

	// Calculate means.
	mongoMean := stat.Mean(mongoStats, nil)
	graphMean := stat.Mean(graphStats, nil)
	comp := Comparison{
		Mongo: mongoMean,
		Graph: graphMean,
	}
	results.UserAssets = comp
	return nil
}

// assetCommentsAuthors gathers stats on the single asset view requests.
func assetCommentsAuthors() error {

	// Gather statistics.
	var mongoStats []float64
	var graphStats []float64
	for i := 1; i <= numRequests; i++ {
		asset, err := getAsset()
		if err != nil {
			return errors.Wrap(err, "Could not retrieve a random asset")
		}
		graphurl := fmt.Sprintf(assetGraphURL, shelfHost, asset)
		mongourl := fmt.Sprintf(assetMongoURL, shelfHost, asset)
		graphStat, err := timeResponse(graphurl, false)
		if err != nil {
			return errors.Wrap(err, "Could not time graph request")
		}
		mongoStat, err := timeResponse(mongourl, false)
		if err != nil {
			return errors.Wrap(err, "Could not time mongorequest")
		}
		mongoStats = append(mongoStats, mongoStat)
		graphStats = append(graphStats, graphStat)
	}

	// Calculate means.
	mongoMean := stat.Mean(mongoStats, nil)
	graphMean := stat.Mean(graphStats, nil)
	comp := Comparison{
		Mongo: mongoMean,
		Graph: graphMean,
	}
	results.AssetUsersComments = comp
	return nil
}

// parentedComments gathers stats on parented comments view requests.
func parentedComments() error {

	// Gather statistics.
	var mongoStats []float64
	var graphStats []float64
	for i := 1; i <= numRequests; i++ {
		comment, err := getComment()
		if err != nil {
			return errors.Wrap(err, "Could not retrieve a random comment")
		}
		graphurl := fmt.Sprintf(parentedGraphURL, shelfHost, comment)
		mongourl := fmt.Sprintf(parentedMongoURL, shelfHost, comment)
		graphStat, err := timeResponse(graphurl, true)
		if err != nil && graphStat > 0.0 {
			return errors.Wrap(err, "Could not time graph request")
		}
		if err != nil && graphStat <= 0.0 {
			continue
		}
		mongoStat, err := timeResponse(mongourl, true)
		if err != nil && mongoStat > 0.0 {
			return errors.Wrap(err, "Could not time mongorequest")
		}
		if err != nil && mongoStat <= 0.0 {
			continue
		}
		mongoStats = append(mongoStats, mongoStat)
		graphStats = append(graphStats, graphStat)
	}

	// Try again if we didn't find any grandparented comments.
	if len(graphStats) == 0 {
		if timeOut >= 3 {
			return fmt.Errorf("Could not find parents")
		}
		log.Println("Sorry for the delay, still looking for random parents")
		timeOut++
		if err := parentedComments(); err != nil {
			log.Fatal(err)
		}
	}

	// Calculate means.
	mongoMean := stat.Mean(mongoStats, nil)
	graphMean := stat.Mean(graphStats, nil)
	comp := Comparison{
		Mongo: mongoMean,
		Graph: graphMean,
	}
	results.Parents = comp
	return nil
}

// grandParentedComments gathers stats on parented comments view requests.
func grandParentedComments() error {

	// Gather statistics.
	var mongoStats []float64
	var graphStats []float64
	for i := 1; i <= numRequests; i++ {
		comment, err := getComment()
		if err != nil {
			return errors.Wrap(err, "Could not retrieve a random comment")
		}
		graphurl := fmt.Sprintf(grandParentedGraphURL, shelfHost, comment)
		mongourl := fmt.Sprintf(grandParentedMongoURL, shelfHost, comment)
		graphStat, err := timeResponse(graphurl, true)
		if err != nil && graphStat > 0.0 {
			return errors.Wrap(err, "Could not time graph request")
		}
		if err != nil && graphStat <= 0.0 {
			continue
		}
		mongoStat, err := timeResponse(mongourl, true)
		if err != nil && mongoStat > 0.0 {
			return errors.Wrap(err, "Could not time mongorequest")
		}
		if err != nil && mongoStat <= 0.0 {
			continue
		}
		mongoStats = append(mongoStats, mongoStat)
		graphStats = append(graphStats, graphStat)
	}

	// Try again if we didn't find any grandparented comments.
	if len(graphStats) == 0 {
		if timeOut >= 3 {
			return fmt.Errorf("Could not find grandparents")
		}
		log.Println("Sorry for the delay, still looking for random grandparents")
		timeOut++
		if err := grandParentedComments(); err != nil {
			log.Fatal(err)
		}
	}

	// Calculate means.
	mongoMean := stat.Mean(mongoStats, nil)
	graphMean := stat.Mean(graphStats, nil)
	comp := Comparison{
		Mongo: mongoMean,
		Graph: graphMean,
	}
	results.GrandParents = comp
	return nil
}

// grGrandParentedComments gathers stats on parented comments view requests.
func grGrandParentedComments() error {

	// Gather statistics.
	var mongoStats []float64
	var graphStats []float64
	for i := 1; i <= numRequests; i++ {
		comment, err := getComment()
		if err != nil {
			return errors.Wrap(err, "Could not retrieve a random comment")
		}
		graphurl := fmt.Sprintf(grGrandParentedGraphURL, shelfHost, comment)
		mongourl := fmt.Sprintf(grGrandParentedMongoURL, shelfHost, comment)
		graphStat, err := timeResponse(graphurl, true)
		if err != nil && graphStat > 0.0 {
			return errors.Wrap(err, "Could not time graph request")
		}
		if err != nil && graphStat <= 0.0 {
			continue
		}
		mongoStat, err := timeResponse(mongourl, true)
		if err != nil && mongoStat > 0.0 {
			return errors.Wrap(err, "Could not time mongorequest")
		}
		if err != nil && mongoStat <= 0.0 {
			continue
		}
		mongoStats = append(mongoStats, mongoStat)
		graphStats = append(graphStats, graphStat)
	}

	// Try again if we didn't find any grandparented comments.
	if len(graphStats) == 0 {
		if timeOut >= 3 {
			return fmt.Errorf("Could not find grandparents")
		}
		log.Println("Sorry for the delay, still looking for random grandparents")
		timeOut++
		if err := grandParentedComments(); err != nil {
			log.Fatal(err)
		}
	}

	// Calculate means.
	mongoMean := stat.Mean(mongoStats, nil)
	graphMean := stat.Mean(graphStats, nil)
	comp := Comparison{
		Mongo: mongoMean,
		Graph: graphMean,
	}
	results.GrGrandParents = comp
	return nil
}

// printResults formats and output the statistics to standard out.
func printResults() {

	// Print the parameters of the statistics run.
	fmt.Printf("NUMBER OF DOCUMENTS: %d\n", numDocs)
	fmt.Printf("NUMBER OF REQUESTS PER QUERY: %d\n", numRequests)
	fmt.Println("")

	// Print the stats for query 1.
	fmt.Println("=============================================")
	fmt.Println("QUERY #1: All comments authored by a random user")
	fmt.Println("NUMBER OF RELATIONSHIP LEVELS: 1")
	fmt.Println("")
	fmt.Printf("MONGO MEAN RESPONSE TIME (seconds): %.4f\n", results.UserComments.Mongo)
	fmt.Printf("GRAPH MEAN RESPONSE TIME (seconds): %.4f\n", results.UserComments.Graph)
	fmt.Println("")

	// Print the stats for query 2.
	fmt.Println("=============================================")
	fmt.Println("QUERY #2: All assets commented on by a user")
	fmt.Println("NUMBER OF RELATIONSHIP LEVELS: 1")
	fmt.Println("")
	fmt.Printf("MONGO MEAN RESPONSE TIME (seconds): %.4f\n", results.UserAssets.Mongo)
	fmt.Printf("GRAPH MEAN RESPONSE TIME (seconds): %.4f\n", results.UserAssets.Graph)
	fmt.Println("")

	// Print the stats for query 3.
	fmt.Println("=============================================")
	fmt.Println("QUERY #3: All comments on an asset along with corresponding authors")
	fmt.Println("NUMBER OF RELATIONSHIP LEVELS: 1-2")
	fmt.Println("")
	fmt.Printf("MONGO MEAN RESPONSE TIME (seconds): %.4f\n", results.AssetUsersComments.Mongo)
	fmt.Printf("GRAPH MEAN RESPONSE TIME (seconds): %.4f\n", results.AssetUsersComments.Graph)
	fmt.Println("")

	// Print the stats for query 4.
	fmt.Println("=============================================")
	fmt.Println("QUERY #4: All comments parented by a set of comments, the set of comments being authored by the author of the parent of a given comment")
	fmt.Println("NUMBER OF RELATIONSHIP LEVELS: 4")
	fmt.Println("")
	fmt.Printf("MONGO MEAN RESPONSE TIME (seconds): %.4f\n", results.Parents.Mongo)
	fmt.Printf("GRAPH MEAN RESPONSE TIME (seconds): %.4f\n", results.Parents.Graph)
	fmt.Println("")

	// Print the stats for query 5.
	fmt.Println("=============================================")
	fmt.Println("QUERY #5: All comments grandparented by a set of comments, the set of comments being authored by the author of the parent of a given comment")
	fmt.Println("NUMBER OF RELATIONSHIP LEVELS: 5")
	fmt.Println("")
	fmt.Printf("MONGO MEAN RESPONSE TIME (seconds): %.4f\n", results.GrandParents.Mongo)
	fmt.Printf("GRAPH MEAN RESPONSE TIME (seconds): %.4f\n", results.GrandParents.Graph)
	fmt.Println("")

	// Print the stats for query 6.
	fmt.Println("=============================================")
	fmt.Println("QUERY #6: All comments great-grandparented by a set of comments, the set of comments being authored by the author of the parent of a given comment")
	fmt.Println("NUMBER OF RELATIONSHIP LEVELS: 6")
	fmt.Println("")
	fmt.Printf("MONGO MEAN RESPONSE TIME (seconds): %.4f\n", results.GrGrandParents.Mongo)
	fmt.Printf("GRAPH MEAN RESPONSE TIME (seconds): %.4f\n", results.GrGrandParents.Graph)
	fmt.Println("")

}

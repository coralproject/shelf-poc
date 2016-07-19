package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// sendToSponge sends a generated item to sponge for processing.
func sendToSponge(payload []byte, typeIn string) error {

	url := fmt.Sprintf("http://%s/1.0/item/coral_%s", spongeHost, typeIn)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrap(err, "Could not create http request")
	}
	req.Header.Add("content-type", "application/json")
	req.Close = true

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "Could not execute POST request to sponge")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Wrap(err, "Unexpected sponge response")
	}
	return nil
}

// worker processes a job on the jobs channel and publishes the result to the
// results channel.
func worker(jobsIn <-chan Job, resultsIn chan<- error) {
	for j := range jobsIn {
		err := sendToSponge(j.Data, j.Type)
		if err == io.EOF {
			time.Sleep(1 * time.Second)
			err := sendToSponge(j.Data, j.Type)
			if err == io.EOF {
				time.Sleep(10 * time.Second)
				err := sendToSponge(j.Data, j.Type)
				if err == io.EOF {
					resultsIn <- err
				}
			}
			continue
		}
		resultsIn <- err
		time.Sleep(10 * time.Millisecond)
	}
}

// handleErrors listens for errors generated in uploading items.
func handleErrors() {
	go func() {
		for err := range results {
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}

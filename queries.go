package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	// "encoding/json"
)

func json_query(endpoint string) (string, error) {
	log.SetPrefix("json_query: ")
	log.SetFlags(0)

	//Empty string check
	if endpoint == "" {
		return "", errors.New("Empty endpoint")
	}

	//Send request to endpoint to get a response
	log.Printf("Request to %v", endpoint)
	response, err := http.Get(endpoint)

	//Error check
	if err != nil {
		return "", err
	}

	//Read the response's data
	respData, err := io.ReadAll(response.Body)

	//Error check
	if err != nil {
		return "", err
	}

	return string(respData), nil
}

func poke_query(id string) (string, error) {
	endpoint := SITE + "pokemon/" + id
	return json_query(endpoint)
}

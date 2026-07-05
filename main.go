package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func json_query(endpoint string)(string, error){
	//Empty string check
	if (endpoint == "") {
		return "", errors.New("Empty endpoint")
	}

	//Send request to endpoint to get a response
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

func main() {
	log.SetPrefix("pk7db: ")
	//Format: Ldate, Ltime, Lshortfile
	log.SetFlags(0b11001000)

	var endpoint string

	endpoint = "https://pokeapi.co/api/v2/pokemon/umbreon"

	json, err := json_query(endpoint)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Print(json)
}

/* important commands:
go build: creates executable
go get (external link): imports external package and updates go.mod and go.sum
go test -v ./...: test there 
*/

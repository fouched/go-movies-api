package main

import (
	"log"
	"net/http"
)

const port = ":9080"

type application struct {
	Domain string
}

func main() {

	// set application config
	var app application

	// read from command line

	// connect to database

	app.Domain = "example.com"

	// start web server
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}

}

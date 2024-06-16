package main

import (
	"database/sql"
	"flag"
	"github.com/fouched/go-movies-api/internal/repo"
	"log"
	"net/http"
)

const port = ":9080"

type application struct {
	DSN    string
	Domain string
}

func main() {
	// set application config
	var app application

	dbPool, err := initApp(app)
	if err != nil {
		log.Fatal(err)
	}
	// we have database connectivity, close it after app stops
	defer dbPool.Close()

	app.Domain = "example.com"
	log.Println("Starting application on", port)
	// start web server
	err = http.ListenAndServe(port, app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

func initApp(app application) (*sql.DB, error) {
	// read from command line, var to be populated, the flag, default value and some help text
	flag.StringVar(&app.DSN,
		"dsn",
		"host=localhost port=5432 user=fouche password=javac dbname=movies sslmode=disable timezone=UTC connect_timeout=5",
		"Database connection string")
	flag.Parse()

	dbPool, err := repo.CreateDbPool(app.DSN)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	} else {
		log.Println("Connected to database!")
	}
	return dbPool, err
}

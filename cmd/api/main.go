package main

import (
	"database/sql"
	"flag"
	"github.com/fouched/go-movies-api/internal/repo"
	"log"
	"net/http"
	"time"
)

const port = ":9080"

type application struct {
	DSN          string
	Domain       string
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

var app application

func main() {
	dbPool, err := initApp()
	if err != nil {
		log.Fatal(err)
	}
	// we have database connectivity, close it after app stops
	defer dbPool.Close()

	log.Println("Starting application on", port)
	// start web server
	err = http.ListenAndServe(port, app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

func initApp() (*sql.DB, error) {
	// read from command line, var to be populated, the flag, default value and some help text
	flag.StringVar(&app.DSN,
		"dsn",
		"host=localhost port=5432 user=fouche password=javac dbname=movies sslmode=disable timezone=UTC connect_timeout=5",
		"Database connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issues")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.Parse()

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "__Host-refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	dbPool, err := repo.CreateDbPool(app.DSN)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	} else {
		log.Println("Connected to database!")
	}
	return dbPool, err
}

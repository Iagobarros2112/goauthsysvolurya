package main

//Connect to PostgreSQL database and Redis store
//Having structured our application, it's time
//to connect to our database and a redis store.
//Redis is needed for the temporary storage of tokens
//and cookies. That is surely more performant than storing
//them in the database.

//To begin with, let's create a config type in main.go.
//This custom type will be made available to all
//our routes via another type called application by
//binding the routes as functions to the type.
//One of Go's paradigms for OOP (Object-oriented Programming):

import (
	"os"

	"goauthbackend/internal/jsonlog"

	"github.com/redis/go-redis/v9"
)

const version = "1.0.0"

// `config` type to house all our app's configurations

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	redisURL string
}

// Main `application` type
type application struct {
	config      config
	logger      *jsonlog.Logger
	redisClient *redis.Client
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	cfg, err := updateConfigWithEnvVariables()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	db, err := openDB(*cfg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	opt, err := redis.ParseURL(cfg.redisURL)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	client := redis.NewClient(opt)

	logger.PrintInfo("redis connection pool established", nil)

	app := &application{
		config:      *cfg,
		logger:      logger,
		redisClient: client,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

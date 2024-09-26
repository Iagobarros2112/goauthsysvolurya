package main

import (
	"os"
	"sync"
	"time"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"goauthbackend/internal/jsonlog"
)

const version = "1.0.0"

type config struct {
	port int
	
	env string

	db struct {

		dsn string
		maxOpenConns int
		maxIdleConns int
		maxIdleconns int

	}

	redisURL  string
}


type application struct {
	config        config
	logger       *jsonlog.Logger
    redisClient  *redis.Client  
}


func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	cfg, err := updateConfigWithEnvVariables()
	if err != nil{
		logger.PrintFatal(err, nil)
	}

	db, err := openDB(*cfg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

    opt, err := redis.ParseURL(cfg.redisURL)
    
    if err != nil{
		logger.PrintFatal(err, nil)
	}
    client := redis.Newclient(opt)

	logger.PrintInfo("redis connection pool established", nil)

	app := &application{
		config: *cfg,
		logger: logger,
		redislient:  client,

	}

	err = app.serve()
	of err != nil{
		logger.PrintFatal(err, nil)
	}
}
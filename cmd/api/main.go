package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Blue-Davinci/SocialAid/internal/logger"
	"go.uber.org/zap"
)

type config struct {
	port int
	env  string
	api  struct {
		name   string
		author string
	}
	db struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	encryption struct {
		key string
	}
}

type application struct {
	config config
	logger *zap.Logger
}

func main() {
	// set up our logger
	logger, err := logger.InitJSONLogger()
	if err != nil {
		fmt.Println("Error initializing logger")
		return
	}
	// set up our config
	var cfg config
	// Port & env
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// API configuration
	flag.StringVar(&cfg.api.name, "api-name", "OptiVest", "API name")
	flag.StringVar(&cfg.api.author, "api-author", "Brian Karicha", "API author")
	// Database configuration
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("SOCIALAID_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	// Encryption key
	flag.StringVar(&cfg.encryption.key, "encryption-key", os.Getenv("SOCIALAID_DATA_ENCRYPTION_KEY"), "Encryption key")

	// create dependancies
	app := &application{
		config: cfg,
		logger: logger,
	}
	fmt.Println("PlaceHolder for the main application")
	app.logger.Info("Starting the application", zap.String("env", app.config.env), zap.Int("port", app.config.port))
}

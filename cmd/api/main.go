package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Blue-Davinci/SocialAid/internal/data"
	"github.com/Blue-Davinci/SocialAid/internal/database"
	"github.com/Blue-Davinci/SocialAid/internal/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *zap.Logger
	models data.Models
}

func main() {
	// set up our logger
	logger, err := logger.InitJSONLogger()
	if err != nil {
		fmt.Println("Error initializing logger")
		return
	}
	// Load the environment variables from the .env file
	getCurrentPath(logger)
	// set up our config
	var cfg config
	// Port & env
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// API configuration
	flag.StringVar(&cfg.api.name, "api-name", "SocialAid", "API name")
	flag.StringVar(&cfg.api.author, "api-author", "Brian Karicha", "API author")
	// Database configuration
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("SOCIALAID_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	// Encryption key
	flag.StringVar(&cfg.encryption.key, "encryption-key", os.Getenv("SOCIALAID_DATA_ENCRYPTION_KEY"), "Encryption key")
	// CORS configuration
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil

	})

	// create our connection pull
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("dsn", cfg.db.dsn))
	}
	logger.Info("database connection pool established", zap.String("dsn", cfg.db.dsn))
	// create dependancies
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}
	// start the server
	err = app.server()
	if err != nil {
		logger.Fatal("Error while starting server.", zap.String("error", err.Error()))
	}
	app.logger.Info("Starting the application", zap.String("env", app.config.env), zap.Int("port", app.config.port))
}

func openDB(cfg config) (*database.Queries, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	// Use ping to establish new conncetions
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	queries := database.New(db)
	return queries, nil
}

// getCurrentPath invokes getEnvPath to get the path to the .env file based on the current working directory.
// After that it loads the .env file using godotenv.Load to be used by the initFlags() function
func getCurrentPath(logger *zap.Logger) string {
	currentpath := getEnvPath(logger)
	if currentpath != "" {
		err := godotenv.Load(currentpath)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("path", currentpath))
		}
	} else {

		logger.Error("Path Error", zap.String("path", currentpath), zap.String("error", "unable to load .env file"))
	}
	logger.Info("Loading Environment Variables", zap.String("path", currentpath))
	return currentpath
}

// getEnvPath returns the path to the .env file based on the current working directory.
func getEnvPath(logger *zap.Logger) string {
	dir, err := os.Getwd()
	if err != nil {
		logger.Fatal(err.Error(), zap.String("path", dir))
		return ""
	}
	if strings.Contains(dir, "cmd/api") || strings.Contains(dir, "cmd") {
		return ".env"
	}
	return filepath.Join("cmd", "api", ".env")
}

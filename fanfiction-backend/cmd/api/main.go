package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/data"
	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/jsonlog"
	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/mailer"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

// config is a struct containing all the command line variables used to configure the api
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// application is a struct used to manage all application-wide dependencies to make them available to the handlers and other functions
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	var cfg config

	webPort := os.Getenv("PORT")
	cfg.port, _ = strconv.Atoi(webPort)
	if cfg.port == 0 {
		cfg.port = 4000
	}

	env := os.Getenv("ENV")
	if env == "" {
		cfg.env = "development"
	}

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 10, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 100, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "1h", "PostgreSQL max connection idle time")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	logger.PrintInfo("Database connection pool established", nil)

	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	sender := "FFWriter <no-reply@ffwriter.net>"

	// ! Mailtrap testing
	// host = "sandbox.smtp.mailtrap.io"
	// port = 2525
	// username = "5c5e586e13d508"
	// password = "a196a1fa87d620"

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(host, port, username, password, sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := "host=" + host +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbName +
		" port=" + port +
		" sslmode=disable"

	if cfg.env == "development" {
		dsn = os.Getenv("FFWRITER_DB_DSN") // For local purposes
	}
	db, err := sql.Open("postgres", dsn)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

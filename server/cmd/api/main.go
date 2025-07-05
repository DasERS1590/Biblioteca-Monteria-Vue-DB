package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"biblioteca/internal/auth"
	"biblioteca/internal/data"
	"biblioteca/internal/jsonlog"
	"biblioteca/internal/mailer"
	"biblioteca/internal/vcs"
    _ "biblioteca/docs"
	_ "github.com/go-sql-driver/mysql"
)

var (
	version = vcs.Version()
)

type config struct {
	port int
	env  string
	db   struct {
		dns          string
		maxIdleConns int
		maxOpenConns int
		maxIdleTime  string
	}

	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}

	jwt struct {
		secret string
		exp    time.Duration
		iss    string
	}
}

type application struct {
	config        config
	logger        *jsonlog.Logger
	models        data.Models
	mailer        mailer.Mailer
	wg            sync.WaitGroup
	authenticator auth.Authenticator
}


// @title BibliotecaAPI
// @version 1.0
// @description Rest-Api for library
// @BasePath /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by your JWT token.
func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "HTTP network address port for API")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dns, "dns", "root:admin@(127.0.0.1:3306)/biblioteca?parseTime=true", "MySQL data source name")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	//flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	//flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	//flag.StringVar(&cfg.smtp.username, "smtp-username", "00000000000000", "SMTP username")
	//flag.StringVar(&cfg.smtp.password, "smtp-password", "00000000000000", "SMTP password")
	//flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <example@example.com>", "SMTP sender")

	flag.StringVar(&cfg.jwt.secret, "jwt", "secreto123", "JWT secrect")
	flag.DurationVar(&cfg.jwt.exp, "exp", time.Hour, "Duración del token (ej: 1h)")
	flag.StringVar(&cfg.jwt.iss, "iss", "miApp", "Emisor del token")

	cfg.cors.trustedOrigins = []string{"http://localhost:3000"} 

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.jwt.secret,
		cfg.jwt.iss,
		cfg.jwt.iss,
	)

	app := &application{
		config:        cfg,
		logger:        logger,
		models:        data.NewModels(db),
		mailer:        mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		authenticator: jwtAuthenticator,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.db.dns)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}

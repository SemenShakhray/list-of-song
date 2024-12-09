package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/goloop/env"
)

type Config struct {
	DB        DB
	API       API
	Server    Server
	Migration Migration
}

type DB struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type API struct {
	Call bool
	Host string
	Port string
}

type Server struct {
	Host        string
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type Migration struct {
	Dir  string
	DSN  string
	Name string
}

func Load() Config {
	cfg := Config{}
	err := env.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}
	cfg = Config{
		DB: DB{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
		},
		API: API{
			Host: os.Getenv("API_HOST"),
			Port: os.Getenv("API_PORT"),
		},
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Migration: Migration{
			Dir:  os.Getenv("MIGRATION_DIR"),
			DSN:  os.Getenv("MIGRATION_DSN"),
			Name: os.Getenv("MIGRATION_NAME"),
		},
	}

	callApi, err := strconv.ParseBool(os.Getenv("API_CALL"))
	if err != nil {
		log.Fatal(err)
	}
	timeout, err := time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	cfg.API.Call = callApi
	cfg.Server.Timeout = timeout
	cfg.Server.IdleTimeout = idleTimeout

	log.Println(cfg)
	return cfg
}

package config

import (
	"log"
	"os"
	"time"

	"github.com/goloop/env"
)

type Config struct {
	DB     DB
	API    API
	Server Server
}

type DB struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type API struct {
	Host string
	Port string
}

type Server struct {
	Host        string
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func Load() Config {
	cfg := Config{}
	err := env.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}
	cfg = Config{
		DB: DB{
			Host: env.Get("DB_HOST"),
			Port: env.Get("DB_PORT"),
			User: env.Get("DB_USER"),
			Pass: env.Get("DB_PASS"),
			Name: env.Get("DB_NAME"),
		},
		API: API{
			Host: env.Get("API_HOST"),
			Port: env.Get("API_PORT"),
		},
		Server: Server{
			Host: env.Get("Server_HOST"),
			Port: env.Get("Server_PORT"),
		},
	}

	timeout, err := time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	idleTimeout, err := time.ParseDuration(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}
	cfg.Server.Timeout = timeout
	cfg.Server.IdleTimeout = idleTimeout

	return cfg
}

package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/SemenShakhray/list-of-song/internal/app"
)

// @title			Songs Library API
// @version		1.0.0
// @description	API for managing a song library
// @contact.url	https://github.com/SemenShakhray
// @BasePath		/
func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed created app: %v", err)
	}

	go func() {
		err := app.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed ListenAndServe: %v", err)
		}
	}()

	sig := <-app.Sigint
	log.Printf("Received signal: %v", sig)

	err = app.Stop()
	log.Printf("HTTP server shutdown error: %v", err)

}

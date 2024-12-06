package main

import (
	"errors"
	"listsongs/internal/app"
	"log"
	"net/http"
)

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

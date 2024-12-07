package router

import (
	"net/http"

	_ "github.com/SemenShakhray/list-of-song/docs"
	"github.com/SemenShakhray/list-of-song/internal/api/handlers"
	"github.com/SemenShakhray/list-of-song/internal/api/middleware"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(h *handlers.Handler) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.WithLogging(h.Log))

	r.Put("/songs/{id}", http.HandlerFunc(h.Update))
	r.Get("/songs", http.HandlerFunc(h.GetAll))
	r.Get("/songs/{song}", http.HandlerFunc(h.GetText))
	r.Post("/songs", http.HandlerFunc(h.AddSong))
	r.Delete("/songs/{id}", http.HandlerFunc(h.Delete))

	r.Get("/swagger", httpSwagger.WrapHandler)

	return r
}

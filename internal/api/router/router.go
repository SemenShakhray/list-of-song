package router

import (
	"listsongs/internal/api/handlers"
	"listsongs/internal/api/middleware"
	"net/http"

	"github.com/go-chi/chi"
)

func NewRouter(h *handlers.Handler) *chi.Mux {

	r := chi.NewRouter()

	r.Put("/songs/{id}", middleware.WithLogging(h.Log, http.HandlerFunc(h.Update)))
	r.Get("/songs", middleware.WithLogging(h.Log, http.HandlerFunc(h.GetAll)))
	r.Get("/songs/{song}", middleware.WithLogging(h.Log, http.HandlerFunc(h.GetText)))
	r.Post("/songs", middleware.WithLogging(h.Log, http.HandlerFunc(h.AddSong)))
	r.Delete("/songs/{id}", middleware.WithLogging(h.Log, http.HandlerFunc(h.Delete)))

	return r
}

package router

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	event "github.com/Dyastin-0/gitsense/internal/event"
	"github.com/Dyastin-0/gitsense/internal/middleware"
)

func Output(client *mongo.Client) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.JWT)
	router.Get("/", event.Get(client))
	return router
}

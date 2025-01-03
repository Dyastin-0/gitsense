package router

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	middleware "github.com/Dyastin-0/gitsense/internal/middleware"
	callback "github.com/Dyastin-0/gitsense/internal/webhook/callback"
)

func Callback(client *mongo.Client) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.JWT)

	router.Post("/{username}/{repository}/hooks/{name}", callback.Handler(client))
	return router
}

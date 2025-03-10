package router

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/oauth2"

	middleware "github.com/Dyastin-0/gitsense/internal/middleware"
	repository "github.com/Dyastin-0/gitsense/internal/repository"
)

func Repository(config *oauth2.Config, client *mongo.Client) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.JWT)
	router.Get("/", repository.Get(config, client))
	return router
}

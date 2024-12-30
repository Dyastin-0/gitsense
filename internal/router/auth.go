package router

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/oauth2"

	auth "gitsense/internal/auth"
)

func Auth(config *oauth2.Config, client *mongo.Client) chi.Router {
	router := chi.NewRouter()

	router.Post("/", auth.Handler(config))
	router.Get("/callback", auth.Callback(config, client))

	return router
}

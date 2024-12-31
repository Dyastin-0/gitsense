package router

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/oauth2"

	auth "gitsense/internal/auth"
)

func Auth(config *oauth2.Config, client *mongo.Client) chi.Router {
	router := chi.NewRouter()

	router.Post("/refresh", auth.Refresh(client))
	router.Get("/github", auth.Handler(config))
	router.Get("/github/callback", auth.Callback(config, client))

	return router
}

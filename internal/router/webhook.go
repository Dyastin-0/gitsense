package router

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	"github.com/Dyastin-0/gitsense/internal/middleware"
	webhook "github.com/Dyastin-0/gitsense/internal/webhook"
)

func Webhook(config *oauth2.Config, mongoClient *mongo.Client) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.JWT)
	router.Post("/", webhook.Create(config, mongoClient))

	return router
}

package router

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Dyastin-0/gitsense/internal/middleware"
	output "github.com/Dyastin-0/gitsense/internal/output"
)

func Output(client *mongo.Client) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.JWT)
	router.Get("/", output.Get(client))
	return router
}

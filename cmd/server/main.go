package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	router "github.com/Dyastin-0/gitsense/internal/router"

	middleware "github.com/Dyastin-0/gitsense/internal/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	githubOAuthConfig := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes: []string{
			"user",
			"repo",
			"admin:repo_hook",
		},
		Endpoint: github.Endpoint,
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic(err)
	}

	mainRouter := chi.NewRouter()

	mainRouter.Use(middleware.Logger)
	mainRouter.Use(middleware.Credential)
	mainRouter.Use(render.SetContentType(render.ContentTypeJSON))

	mainRouter.Mount("/auth", router.Auth(githubOAuthConfig, client))
	mainRouter.Mount("/repository", router.Repository(githubOAuthConfig, client))
	mainRouter.Mount("/webhook", router.Webhook(githubOAuthConfig, client))
	mainRouter.Mount("/callback", router.Callback(client))
	mainRouter.Mount("/event", router.Event(client))

	port := os.Getenv("PORT")
	if err := http.ListenAndServe(":"+port, mainRouter); err != nil {
		log.Fatal("Server failed: ", err)
	}
}

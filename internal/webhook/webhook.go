package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	types "github.com/Dyastin-0/gitsense/internal/auth/type"
	"github.com/Dyastin-0/gitsense/pkg/util/aes"
	"github.com/Dyastin-0/gitsense/pkg/util/hash"
)

func Create(config *oauth2.Config, mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*types.Claims)
		if !ok || claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user := claims.User

		accessToken := user.GithubAccessToken
		token := &oauth2.Token{
			AccessToken: accessToken,
		}

		req := RequestBody{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		collection := mongoClient.Database("test").Collection("webhooks")
		filter := bson.M{"name": req.Name, "owner": user.Login}
		err = collection.FindOne(r.Context(), filter).Decode(&Webhook{})

		if err == nil {
			http.Error(w, "Webhook already exists", http.StatusConflict)
			return
		}

		webhookPayloadURL := fmt.Sprintf("%s/%s/%s/hooks/%s", os.Getenv("BASE_SERVER_URL"), user.Login, req.Repository, req.Name)

		webhookPayload := WebhookPayload{
			Name:   "web",
			Active: true,
			Events: []string{"push"},
			Config: WebhookConfig{
				URL:         webhookPayloadURL,
				ContentType: "json",
				Secret:      req.Secret,
				InsecureSSL: "0",
			},
		}

		payloadBytes, err := json.Marshal(webhookPayload)
		if err != nil {
			http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
			return
		}

		client := config.Client(r.Context(), token)

		webhookURL := fmt.Sprintf("%s/repos/%s/%s/hooks", os.Getenv("GITHUB_API_URL"), user.Login, req.Repository)
		resp, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			http.Error(w, "Failed to create webhook", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			http.Error(w, fmt.Sprintf("GitHub API error: %s", resp.Status), http.StatusInternalServerError)
			return
		}

		req.SSH.PrivateKey, err = aes.Encrypt(req.SSH.PrivateKey, os.Getenv("ENCRYPTION_KEY"))
		if err != nil {
			http.Error(w, "Failed to encrypt private key", http.StatusInternalServerError)
			return
		}

		req.Secret, err = hash.Generate(req.Secret)
		if err != nil {
			http.Error(w, "Failed to hash secret", http.StatusInternalServerError)
			return
		}

		webhook := Webhook{
			Name:           req.Name,
			Repository:     req.Repository,
			Secret:         req.Secret,
			Owner:          user.Login,
			SSHconfig:      req.SSH,
			WebhookPayload: webhookPayload,
			Created:        time.Now().Unix(),
		}

		_, err = collection.InsertOne(r.Context(), webhook)
		if err != nil {
			http.Error(w, "Failed to insert webhook", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

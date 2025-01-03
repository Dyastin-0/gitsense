package middleware

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Dyastin-0/gitsense/internal/webhook"
	"github.com/go-chi/chi/v5"
)

func Signature(next http.Handler, mongoClient *mongo.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		repository := chi.URLParam(r, "repository")
		login := chi.URLParam(r, "username")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read request body", http.StatusInternalServerError)
			return
		}

		signatureHeader := r.Header.Get("X-Hub-Signature")
		if signatureHeader == "" {
			http.Error(w, "Missing signature header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(signatureHeader, "=")
		if len(parts) != 2 || parts[0] != "sha1" {
			http.Error(w, "Invalid signature format", http.StatusUnauthorized)
			return
		}

		expectedSignature := parts[1]

		collection := mongoClient.Database("test").Collection("webhooks")
		filter := bson.M{"owner": login, "repository": repository, "name": name}
		webhook := webhook.Webhook{}
		err = collection.FindOne(r.Context(), filter).Decode(&webhook)
		if err != nil {
			http.Error(w, "Webhook not found", http.StatusNotFound)
			return
		}

		mac := hmac.New(sha1.New, []byte(webhook.Secret))
		mac.Write(body)
		calculatedSignature := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(calculatedSignature), []byte(expectedSignature)) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

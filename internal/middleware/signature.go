package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt" // For debugging
	"io"
	"net/http"
	"strings"

	"github.com/Dyastin-0/gitsense/internal/webhook"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

		fmt.Println("Request body:", string(body))

		signatureHeader := r.Header.Get("X-Hub-Signature-256")
		if signatureHeader == "" {
			http.Error(w, "Missing signature header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(signatureHeader, "=")
		if len(parts) != 2 || parts[0] != "sha256" {
			http.Error(w, "Invalid signature format", http.StatusUnauthorized)
			return
		}
		expectedSignature := parts[1]

		collection := mongoClient.Database("test").Collection("webhooks")
		filter := bson.M{"owner": login, "repository": repository, "name": name}
		var webhook webhook.Webhook
		err = collection.FindOne(r.Context(), filter).Decode(&webhook)
		if err != nil {
			http.Error(w, "Webhook not found", http.StatusNotFound)
			return
		}

		fmt.Println("Fetched Secret:", webhook.Secret)

		mac := hmac.New(sha256.New, []byte(webhook.Secret))
		mac.Write(body)
		calculatedSignature := hex.EncodeToString(mac.Sum(nil))

		fmt.Println("Calculated Signature:", calculatedSignature)
		fmt.Println("Expected Signature:", expectedSignature)

		if !hmac.Equal([]byte(calculatedSignature), []byte(expectedSignature)) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package callback

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Dyastin-0/gitsense/pkg/util/aes"
	"github.com/Dyastin-0/gitsense/pkg/util/ssh"
	"github.com/go-chi/chi/v5"

	"github.com/Dyastin-0/gitsense/internal/webhook"
)

func Handler(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		login := chi.URLParam(r, "username")
		repository := chi.URLParam(r, "repository")
		name := chi.URLParam(r, "name")

		fmt.Println(login, repository, name)

		webhook := webhook.Webhook{}

		collection := mongoClient.Database("test").Collection("webhooks")
		filter := bson.M{"owner": login, "repository": repository, "name": name}

		err := collection.FindOne(r.Context(), filter).Decode(&webhook)
		if err != nil {
			http.Error(w, "Webhook not found", http.StatusNotFound)
			return
		}

		decodedPrivateKey, err := aes.Decrypt(webhook.SSHconfig.PrivateKey, os.Getenv("ENCRYPTION_KEY"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stdout, stderr, err := ssh.RunCommand(decodedPrivateKey,
			webhook.SSHconfig.IPadress,
			webhook.SSHconfig.HostKey,
			webhook.SSHconfig.User,
			webhook.CallbackScript,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := Response{
			Stdout: stdout,
			Stderr: stderr,
		}

		json.NewEncoder(w).Encode(response)
	}
}

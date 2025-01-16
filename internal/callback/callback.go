package callback

import (
	"context"
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

		go processWebhook(&Job{
			Owner:      login,
			Repository: repository,
			Name:       name,
		}, mongoClient,
			r.Context(),
		)

		w.WriteHeader(http.StatusAccepted)
	}
}

func processWebhook(job *Job, mongoClient *mongo.Client, ctx context.Context) {
	webhook := webhook.Webhook{}

	collection := mongoClient.Database("test").Collection("webhooks")
	filter := bson.M{"owner": job.Owner, "repository": job.Repository, "name": job.Name}

	err := collection.FindOne(context.Background(), filter).Decode(&webhook)
	if err != nil {
		fmt.Println(fmt.Errorf("error decoding webhook: %v", err))
		return
	}

	decodedPrivateKey, err := aes.Decrypt(webhook.SSHconfig.PrivateKey, os.Getenv("ENCRYPTION_KEY"))
	if err != nil {
		fmt.Println(fmt.Errorf("error decrypting private key: %v", err))
		return
	}

	collection = mongoClient.Database("test").Collection("outputs")

	stdout, stderr, err := ssh.RunCommand(decodedPrivateKey,
		webhook.SSHconfig.IPadress,
		webhook.SSHconfig.HostKey,
		webhook.SSHconfig.User,
		webhook.CallbackScript,
	)
	if err != nil {
		fmt.Println(fmt.Errorf("error running command: %v", err))
		return
	}

	output := Output{
		Stdout:  stdout,
		Stderr:  stderr,
		Webhook: job.Name,
		Owner:   job.Owner,
	}

	_, err = collection.InsertOne(ctx, output)
	if err != nil {
		fmt.Println(fmt.Errorf("error inserting output: %v", err))
		return
	}
}

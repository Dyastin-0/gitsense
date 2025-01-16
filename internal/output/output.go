package output

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	types "github.com/Dyastin-0/gitsense/internal/auth/type"
	"github.com/Dyastin-0/gitsense/internal/callback"
)

func Get(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*types.Claims)
		if !ok || claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user := claims.User
		collection := mongoClient.Database("test").Collection("outputs")
		filter := bson.M{"owner": user.Login}

		cursor, err := collection.Find(r.Context(), filter)
		if err != nil {
			http.Error(w, "Failed to get outputs", http.StatusInternalServerError)
			return
		}

		var outputs []callback.Output
		err = cursor.All(r.Context(), &outputs)
		if err != nil {
			http.Error(w, "Failed to parse outputs", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(outputs)
	}
}

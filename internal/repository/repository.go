package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	types "gitsense/internal/auth/type"
)

func Get(config *oauth2.Config, mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*types.Claims)
		if !ok || claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := claims.User.GithubAccessToken

		token := &oauth2.Token{
			AccessToken: accessToken,
		}

		client := config.Client(r.Context(), token)

		resp, err := client.Get(os.Getenv("GITHUB_API_URL") + "/user/repos")
		if err != nil {
			http.Error(w, "Failed to get repositories", http.StatusInternalServerError)
			fmt.Printf("Error fetching repositories: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("GitHub API error: %s", resp.Status), http.StatusInternalServerError)
			return
		}

		var repos []Repository
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			http.Error(w, "Failed to parse repositories", http.StatusInternalServerError)
			fmt.Println("Error decoding JSON:", err)
			return
		}

		json.NewEncoder(w).Encode(repos)
	}
}

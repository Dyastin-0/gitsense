package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"

	usr "github.com/Dyastin-0/gitsense/internal/type/user"
	tken "github.com/Dyastin-0/gitsense/pkg/util/token"
)

func Handler(config *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func Callback(config *oauth2.Config, mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		token, err := config.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}

		client := config.Client(r.Context(), token)
		response, err := client.Get("https://api.github.com/user")
		if err != nil {
			return
		}
		defer response.Body.Close()

		user := usr.Type{}
		if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
			http.Error(w, "Failed to parse user information", http.StatusInternalServerError)
			return
		}

		user.GithubAccessToken = token.AccessToken

		collection := mongoClient.Database("test").Collection("users")
		filter := bson.M{"github_id": user.GithubID}

		err = collection.FindOne(r.Context(), filter).Err()
		if err == mongo.ErrNoDocuments {
			user.Created = time.Now().Unix()
			_, err := collection.InsertOne(r.Context(), user)
			if err != nil {
				http.Error(w, "Failed to insert new user", http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			http.Error(w, "Failed to find user", http.StatusInternalServerError)
			return
		}

		refreshToken, err := tken.Generate(&user, os.Getenv("REFRESH_TOKEN_KEY"), 24*time.Hour)
		if err != nil {
			http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
			return
		}

		update := bson.M{"$push": bson.M{"refresh_tokens": refreshToken}, "$set": bson.M{"github_access_token": user.GithubAccessToken}}
		_, err = collection.UpdateOne(r.Context(), filter, update)
		if err != nil {
			http.Error(w, "Failed to update user with refresh token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "rt",
			Value:    refreshToken,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			MaxAge:   24 * 60 * 60,
			Domain:   os.Getenv("DOMAIN"),
			Path:     "/",
		})

		redirectURL := os.Getenv("BASE_CLIENT_URL") + "/home"
		http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)

	}
}

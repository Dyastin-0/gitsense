package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID                string `json:"_id" bson:"_id,omitempty"`
	GithubID          int    `json:"github_id" bson:"github_id"`
	GithubAccessToken string `json:"github_access_token" bson:"github_access_token"`
	Login             string `json:"login" bson:"login"`
	Name              string `json:"name" bson:"name"`
	Email             string `json:"email" bson:"email"`
	Avatar            string `json:"avatar_url" bson:"avatar_url"`
}

type Claims struct {
	User User  `json:"user" bson:"user"`
	Exp  int64 `json:"exp" bson:"exp"`
	jwt.RegisteredClaims
}

type Response struct {
	AccessToken string `json:"accessToken"`
	User        User   `json:"user"`
}

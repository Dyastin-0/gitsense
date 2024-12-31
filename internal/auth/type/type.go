package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID       string `bson:"_id,omitempty"`
	GithubID int    `bson:"github_id"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Avatar   string `bson:"avatar_url"`
}

type Claims struct {
	User User  `json:"user"`
	Exp  int64 `json:"exp"`
	jwt.RegisteredClaims
}

type Response struct {
	AccessToken string `json:"accessToken"`
	User        User   `json:"user"`
}

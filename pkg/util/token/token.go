package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	auth "gitsense/internal/auth/type"
	user "gitsense/internal/type/user"
)

func Generate(user *user.Type, secret string, expiration time.Duration) (string, error) {
	claims := auth.Claims{
		User: auth.User{
			Name:     user.Name,
			Email:    user.Email,
			GithubID: user.GithubID,
			ID:       user.ID,
		},
		Exp: time.Now().Add(expiration).Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Filespace",
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

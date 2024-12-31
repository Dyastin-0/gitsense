package user

type Type struct {
	ID           string   `json:"_id,omitempty" bson:"_id,omitempty"`
	Login        string   `json:"login" bson:"login"`
	GithubID     int      `json:"id" bson:"github_id"`
	Name         string   `json:"name" bson:"name"`
	Email        string   `json:"email" bson:"email"`
	Avatar       string   `json:"avatar_url" bson:"avatar_url"`
	Bio          string   `json:"bio" bson:"bio"`
	URL          string   `json:"html_url" bson:"url"`
	RefreshToken []string `json:"refresh_tokens,omitempty" bson:"refresh_tokens,omitempty"`
}

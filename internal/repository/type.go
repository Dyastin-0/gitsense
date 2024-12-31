package repository

type Repository struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
}

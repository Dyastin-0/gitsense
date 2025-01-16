package callback

type Response struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

type Job struct {
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
	Name       string `json:"name"`
}

type Output struct {
	ID        string `json:"_id" bson:"_id,omitempty"`
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	Webhook   string `json:"webhook"`
	Owner     string `json:"owner"`
	Timestamp int64  `json:"timestamp"`
}

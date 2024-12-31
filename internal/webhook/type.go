package webhook

type Webhook struct {
	Name           string         `json:"_id" bson:"_id"`
	Owner          string         `json:"owner" bson:"owner"`
	Created        int64          `json:"created" bson:"created"`
	WebhookPayload WebhookPayload `json:"config"`
}

type WebhookPayload struct {
	Name   string        `json:"name" bson:"name"`
	Active bool          `json:"active" bson:"active"`
	Events []string      `json:"events" bson:"events"`
	Config WebhookConfig `json:"config" bson:"config"`
}

type WebhookConfig struct {
	URL         string `json:"url" bson:"url"`
	ContentType string `json:"content_type" bson:"content_type"`
	Secret      string `json:"secret" bson:"secret"`
	InsecureSSL string `json:"insecure_ssl" bson:"insecure_ssl"`
}

type RequestBody struct {
	Name       string `json:"name" bson:"name"`
	Repository string `json:"repository" bson:"repository"`
	Secret     string `json:"secret" bson:"secret"`
}

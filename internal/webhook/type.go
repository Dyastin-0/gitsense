package webhook

type Webhook struct {
	Name           string         `json:"name" bson:"name"`
	Owner          string         `json:"owner" bson:"owner"`
	Repository     string         `json:"repository" bson:"repository"`
	Secret         string         `json:"secret" bson:"secret"`
	Created        int64          `json:"created" bson:"created"`
	CallbackScript string         `json:"callback_script" bson:"callback_script"`
	SSHconfig      SSHconfig      `json:"ssh" bson:"ssh"`
	WebhookPayload WebhookPayload `json:"config" bson:"config"`
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
	Name           string    `json:"name" bson:"name"`
	Repository     string    `json:"repository" bson:"repository"`
	Secret         string    `json:"secret" bson:"secret"`
	CallbackScript string    `json:"callback_script" bson:"callback_script"`
	SSH            SSHconfig `json:"ssh" bson:"ssh"`
}

type SSHconfig struct {
	IPadress   string `json:"ip_address" bson:"ip_address"`
	User       string `json:"user" bson:"user"`
	HostKey    string `json:"host_key" bson:"host_key"`
	PrivateKey string `json:"private_key" bson:"private_key"`
}

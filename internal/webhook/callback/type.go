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

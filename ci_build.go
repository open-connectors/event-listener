package main

type CiBuildPayload struct {
	Origin          string        `json:"origin"`
	OriginalID      string        `json:"originalID"`
	Name            string        `json:"name"`
	URL             string        `json:"url"`
	CreatedAt       int64         `json:"createdAt"`
	StartedAt       int64         `json:"startedAt"`
	CompletedAt     int64         `json:"completedAt"`
	TriggeredBy     string        `json:"triggeredBy"`
	Status          string        `json:"status"`
	Conclusion      string        `json:"conclusion"`
	RepoURL         string        `json:"repoUrl"`
	Commit          string        `json:"commit"`
	PullRequestUrls []interface{} `json:"pullRequestUrls"`
	IsDeployment    bool          `json:"isDeployment"`
	Stages          []Stage       `json:"stages"`
}

type Job struct {
	StartedAt   int64  `json:"startedAt"`
	CompletedAt int64  `json:"completedAt"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
}

type Stage struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	StartedAt   int64  `json:"startedAt"`
	CompletedAt int64  `json:"completedAt"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
	URL         string `json:"url"`
	Jobs        []Job  `json:"jobs"`
}

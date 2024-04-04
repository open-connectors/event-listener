package main

import "time"

type CiBuildPayload struct {
	Origin          string        `json:"origin"`
	OriginalID      string        `json:"originalID"`
	Name            string        `json:"name"`
	URL             string        `json:"url"`
	CreatedAt       time.Time     `json:"createdAt"`
	StartedAt       time.Time     `json:"startedAt"`
	CompletedAt     time.Time     `json:"completedAt"`
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
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
}

type Stage struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	URL         string    `json:"url"`
	Jobs        []Job     `json:"jobs"`
}

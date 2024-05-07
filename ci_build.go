package main

type CiBuildPayload struct {
	Origin          string        `json:"origin" dynamodbav:"origin"`
	OriginalID      string        `json:"originalID" dynamodbav:"originalId"`
	Name            string        `json:"name" dynamodbav:"name"`
	URL             string        `json:"url" dynamodbav:"url,omitempty"`
	CreatedAt       int64         `json:"createdAt" dynamodbav:"createdAt,omitempty"`
	StartedAt       int64         `json:"startedAt" dynamodbav:"startedAt,omitempty"`
	CompletedAt     int64         `json:"completedAt" dynamodbav:"completedAt,omitempty"`
	TriggeredBy     string        `json:"triggeredBy" dynamodbav:"triggeredBy,omitempty"`
	Status          string        `json:"status" dynamodbav:"status,omitempty"`
	Conclusion      string        `json:"conclusion" dynamodbav:"conclusion,omitempty"`
	RepoURL         string        `json:"repoUrl" dynamodbav:"repoUrl,omitempty"`
	Commit          string        `json:"commit" dynamodbav:"commit,omitempty"`
	PullRequestUrls []interface{} `json:"pullRequestUrls" dynamodbav:"pullrequestUrls,omitempty"`
	IsDeployment    bool          `json:"isDeployment" dynamodbav:"isDeployment,omitempty"`
	Stages          []Stage       `json:"stages" dynamodbav:"stages,omitempty"`
}

type Job struct {
	StartedAt   int64  `json:"startedAt" dynamodbav:"startedAt,omitempty"`
	CompletedAt int64  `json:"completedAt" dynamodbav:"completedAt,omitempty"`
	Name        string `json:"name" dynamodbav:"name,omitempty"`
	Status      string `json:"status" dynamodbav:"status,omitempty"`
	Conclusion  string `json:"conclusion" dynamodbav:"conslusion,omitempty"`
}

type Stage struct {
	ID          string `json:"id" dynamodbav:"id,omitempty"`
	Name        string `json:"name" dynamodbav:"name,omitempty"`
	StartedAt   int64  `json:"startedAt" dynamodbav:"startedAt,omitempty"`
	CompletedAt int64  `json:"completedAt" dynamodbav:"completedAt,omitempty"`
	Status      string `json:"status" dynamodbav:"status,omitempty"`
	Conclusion  string `json:"conclusion" dynamodbav:"conslusion,omitempty"`
	URL         string `json:"url" dynamodbav:"url,omitempty"`
	Jobs        []Job  `json:"jobs" dynamodbav:"jobs,omitempty"`
}

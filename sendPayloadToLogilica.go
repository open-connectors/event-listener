package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadPlanningData(repoId string, payload []CiBuildPayload) {
	postBody, _ := json.Marshal(payload)
	logilicaUrl := fmt.Sprintf("https://logilica.io/api/import/v1/ci_build/%v/create", repoId)
	contentType := "application/json"

	client := &http.Client{}
	req, err := http.NewRequest("POST", logilicaUrl, bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("x-lgca-domain", "redhat")
	req.Header.Add("X-lgca-token", os.Getenv("LOGILICA_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(body)
}

func GetWorkflowRuns() {
	runsUrl := "https://api.github.com/repos/open-connectors/open-connectors/actions/runs"
	contentType := "Accept: application/vnd.github+json"

	client := &http.Client{}
	req, err := http.NewRequest("GET", runsUrl, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("API_TOKEN")))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var workflows WorkflowRuns
	json.Unmarshal(body, &workflows)
	fmt.Println(workflows.TotalCount)
	GetWorkflowMetadata(workflows.WorkflowRuns[0].ID)
}

func GetWorkflowMetadata(id int64) RunMetadata {
	runsUrl := fmt.Sprintf("https://api.github.com/repos/open-connectors/open-connectors/actions/runs/%d", id)
	contentType := "Accept: application/vnd.github+json"

	client := &http.Client{}
	req, err := http.NewRequest("GET", runsUrl, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("API_TOKEN")))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var run RunMetadata
	json.Unmarshal(body, &run)
	return run
}

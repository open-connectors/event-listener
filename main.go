package main

import (
	"context"
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	_ "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/kelseyhightower/envconfig"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"log"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

// TektonCloudEventData type is used to marshal and unmarshal the payload of
// a Tekton cloud event. It can include a TaskRun or a PipelineRun
type TektonCloudEventData struct {
	TaskRun     *v1beta1.TaskRun     `json:"taskRun,omitempty"`
	PipelineRun *v1beta1.PipelineRun `json:"pipelineRun,omitempty"`
	CustomRun   *v1beta1.CustomRun   `json:"customRun,omitempty"`
}

// newTektonCloudEventData returns a new instance of TektonCloudEventData
func newTektonCloudEventData(ctx context.Context, runObject objectWithCondition) (TektonCloudEventData, error) {
	tektonCloudEventData := TektonCloudEventData{}
	switch v := runObject.(type) {
	case *v1beta1.TaskRun:
		tektonCloudEventData.TaskRun = v
	case *v1beta1.PipelineRun:
		tektonCloudEventData.PipelineRun = v
	case *v1.TaskRun:
		v1beta1TaskRun := &v1beta1.TaskRun{}
		if err := v1beta1TaskRun.ConvertFrom(ctx, v); err != nil {
			return TektonCloudEventData{}, err
		}
		tektonCloudEventData.TaskRun = v1beta1TaskRun
	case *v1.PipelineRun:
		v1beta1PipelineRun := &v1beta1.PipelineRun{}
		if err := v1beta1PipelineRun.ConvertFrom(ctx, v); err != nil {
			return TektonCloudEventData{}, err
		}
		tektonCloudEventData.PipelineRun = v1beta1PipelineRun
	case *v1beta1.CustomRun:
		tektonCloudEventData.CustomRun = v
	}
	return tektonCloudEventData, nil
}

func eventReceiver(ctx context.Context, event cloudevents.Event) error {
	fmt.Println("----------------- Reachinhe")
	var object v1.PipelineRun
	if err := json.Unmarshal(event.DataEncoded, &object); err != nil {
		fmt.Println(err)
	}
	fmt.Println(event.String())
	fmt.Println("SPEC:", object.Spec)
	fmt.Println("Status:", object.Status)
	cdevent, _ := newTektonCloudEventData(ctx, event)
	fmt.Println("EVENT", cdevent)
	// PrepareCiBuildData(object)
	fmt.Println("------------------")
	return nil
}

// func PrepareCiBuildData(object v1.PipelineRun) {
// 	payload = CiBuildPayload{
// 	Origin          :"Tekton",
// 	OriginalID      string        `json:"originalID"`
// 	Name            : object.Spec.Name,
// 	URL             string        `json:"url"`
// 	CreatedAt       time.Time     `json:"createdAt"`
// 	StartedAt       time.Time     `json:"startedAt"`
// 	CompletedAt     time.Time     `json:"completedAt"`
// 	TriggeredBy     :,
// 	Status          :,
// 	Conclusion      :,
// 	RepoURL         :,
// 	Commit          string        `json:"commit"`
// 	PullRequestUrls []interface{} `json:"pullRequestUrls"`
// 	IsDeployment    : true,
// 	Stages          []Stage       `json:"stages"`
// 	}
// }

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	log.Print("Starting Event Listener")
	ctx := context.Background()

	p, err := cloudevents.NewHTTP(cloudevents.WithPort(env.Port), cloudevents.WithPath(env.Path))
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}
	c, err := cloudevents.NewClient(p,
		cloudevents.WithUUIDs(),
		cloudevents.WithTimeNow(),
	)
	if err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	log.Printf("listening on :%d%s\n", env.Port, env.Path)
	if err := c.StartReceiver(ctx, eventReceiver); err != nil {
		log.Fatalf("failed to start receiver: %s", err.Error())
	}

	<-ctx.Done()
}

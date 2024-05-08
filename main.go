package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type Data struct {
	Pipelinerun v1.PipelineRun `json:"pipelineRun"`
}

func eventReceiver(ctx context.Context, event cloudevents.Event) error {
	var dat Data
	if err := json.Unmarshal(event.DataEncoded, &dat); err != nil {
		fmt.Println("Ignore")
	}
	var client *dynamodb.Client
	var err error

	if client, err = newclient(); err != nil {
		log.Fatalf("failed to create dynamoclient: %s", err.Error())
	}
	fmt.Println("Pipleine run", dat.Pipelinerun)
	InsertRecordInDatabase(dat.Pipelinerun, client)
	return nil
}

func InsertRecordInDatabase(object v1.PipelineRun, client *dynamodb.Client) {
	// tables, err := listTables(client, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var tableexists bool
	// for _, val := range tables {
	// 	if val == "TektonCI" {
	// 		tableexists = true
	// 	}
	// }
	// if !tableexists {
	// 	input := &dynamodb.CreateTableInput{
	// 		AttributeDefinitions: []types.AttributeDefinition{
	// 			{
	// 				AttributeName: aws.String("Origin"),
	// 				AttributeType: types.ScalarAttributeTypeS,
	// 			},
	// 			{
	// 				AttributeName: aws.String("OriginalID"),
	// 				AttributeType: types.ScalarAttributeTypeS,
	// 			},
	// 			{
	// 				AttributeName: aws.String("Name"),
	// 				AttributeType: types.ScalarAttributeTypeS,
	// 			},
	// 		},
	// 		KeySchema: []types.KeySchemaElement{
	// 			{
	// 				AttributeName: aws.String("OriginalID"),
	// 				KeyType:       types.KeyTypeHash,
	// 			},
	// 			{
	// 				AttributeName: aws.String("Name"),
	// 				KeyType:       types.KeyTypeHash,
	// 			},
	// 		},
	// 		ProvisionedThroughput: &types.ProvisionedThroughput{
	// 			ReadCapacityUnits:  aws.Int64(10),
	// 			WriteCapacityUnits: aws.Int64(10),
	// 		},
	// 		TableName: aws.String("TektonCI"),
	// 	}
	// 	err = createTable(client, "TektonCI", input)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
	item := PrepareCiBuildData(object)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		fmt.Println("failed to marshal Record, %w", err)
	}
	fmt.Println("Inserting in the database")
	fmt.Println("-----------NEW-----------------")
	fmt.Println(av)
	fmt.Println("-----------BETWEEN----------------")
	fmt.Println("Response from put api ", putItem(client, "TektonCI", av))
}

func PrepareCiBuildData(obj v1.PipelineRun) CiBuildPayload {
	payload := CiBuildPayload{
		Origin:          "Tekton",
		OriginalID:      string(obj.UID),
		Name:            obj.Name,
		URL:             obj.Status.Provenance.RefSource.URI,
		CreatedAt:       obj.Status.StartTime.Time.Unix(),
		StartedAt:       obj.Status.StartTime.Time.Unix(),
		CompletedAt:     obj.Status.CompletionTime.Time.Unix(),
		Status:          string(obj.Status.Conditions[0].Type),
		Conclusion:      string(obj.Status.Conditions[0].Status),
		RepoURL:         obj.Status.Provenance.RefSource.URI,
		Commit:          "",
		PullRequestUrls: nil,
		IsDeployment:    true,
	}
	triggeredBy := TriggeredBy{
		Name:         "Pipelines Operator",
		Email:        "dummy@redhat.com",
		AccountId:    "dummy@redhat.com",
		LastActivity: obj.Status.Conditions[0].LastTransitionTime.Inner.Unix(),
	}
	payload.TriggeredBy = triggeredBy

	var tasks []Job
	for _, val := range obj.Status.ChildReferences {
		job := Job{
			StartedAt:   obj.Status.StartTime.Time.Unix(),
			CompletedAt: obj.Status.CompletionTime.Time.Unix(),
			Name:        val.Name,
			Status:      string(obj.Status.Conditions[0].Status),
			Conclusion:  obj.Status.Conditions[0].Reason,
		}
		tasks = append(tasks, job)
	}
	var stg []Stage
	stage := Stage{
		ID:          string(obj.UID),
		Name:        obj.Name,
		StartedAt:   obj.Status.StartTime.Time.Unix(),
		CompletedAt: obj.Status.CompletionTime.Time.Unix(),
		Status:      string(obj.Status.Conditions[0].Status),
		Conclusion:  obj.Status.Conditions[0].Reason,
		URL:         obj.Status.Provenance.RefSource.URI,
		Jobs:        tasks,
	}
	stg = append(stg, stage)
	payload.Stages = stg
	return payload
}

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

	var client *dynamodb.Client

	if client, err = newclient(); err != nil {
		log.Fatalf("failed to create dynamoclient: %s", err.Error())
	}

	go func() {
		t := time.Tick(5 * time.Minute)
		for {
			select {
			case <-t:
				fmt.Println("Logilica Upload")
				LogilicaUpload(client)
			case <-ctx.Done():
				return
			}
		}
	}()

	log.Printf("listening on :%d%s\n", env.Port, env.Path)
	if err := c.StartReceiver(ctx, eventReceiver); err != nil {
		log.Fatalf("failed to start receiver: %s", err.Error())
	}

	<-ctx.Done()
}

func LogilicaUpload(client *dynamodb.Client) {
	payload := getCiBuildPayload(client)
	UploadPlanningData("872a7985dd8a58328dea96015b738c317039fb5a", payload)
}

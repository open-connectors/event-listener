// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
// snippet-start:[dynamodb.go.load_items]
package main

// snippet-start:[dynamodb.go.load_items.imports]
import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"fmt"
)

// snippet-end:[dynamodb.go.load_items.imports]

// snippet-start:[dynamodb.go.load_items.struct]
// Create struct to hold info about new item
type Item struct {
	Year   int     `dynamodbav:"year"`
	Title  string  `dynamodbav:"title"`
	Plot   string  `dynamodbav:"plot,omitempty"`
	Rating float64 `dynamodbav:"rating,omitempty"`
}

// snippet-end:[dynamodb.go.load_items.func]

func newclient() (*dynamodb.Client, error) {
	region := os.Getenv("REGION")
	url := os.Getenv("URL")
	accsKeyID := os.Getenv("ACCESSKEYID")
	secretAccessKey := os.Getenv("SECRETACCESSKEY")
	fmt.Println(region, "REGION")
	fmt.Println(url, "URL")
	fmt.Println(accsKeyID, "ACCESSKEYID")
	fmt.Println(secretAccessKey, "SECRETACCESSKEY")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: url}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: accsKeyID, SecretAccessKey: secretAccessKey, SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	c := dynamodb.NewFromConfig(cfg)
	return c, nil
}

func main() {
	c, err := newclient()
	if err != nil {
		fmt.Println(err)
	}
	// Add each item to Movies table:
	// tableName := "Movies"
	// _, err = c.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
	// 	AttributeDefinitions: []types.AttributeDefinition{{
	// 		AttributeName: aws.String("year"),
	// 		AttributeType: types.ScalarAttributeTypeN,
	// 	}, {
	// 		AttributeName: aws.String("title"),
	// 		AttributeType: types.ScalarAttributeTypeS,
	// 	}},
	// 	KeySchema: []types.KeySchemaElement{{
	// 		AttributeName: aws.String("year"),
	// 		KeyType:       types.KeyTypeHash,
	// 	}, {
	// 		AttributeName: aws.String("title"),
	// 		KeyType:       types.KeyTypeRange,
	// 	}},
	// 	TableName: aws.String(tableName),
	// 	ProvisionedThroughput: &types.ProvisionedThroughput{
	// 		ReadCapacityUnits:  aws.Int64(10),
	// 		WriteCapacityUnits: aws.Int64(10),
	// 	},
	// })
	// if err != nil {
	// 	log.Printf("Couldn't create table %v. Here's why: %v\n", tableName, err)
	// }
	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	// // Create DynamoDB client
	// svc := dynamodb.New(sess)
	// item := Item{
	// 	Year:   2015,
	// 	Title:  "The Big New Movie",
	// 	Plot:   "Nothing happens at all.",
	// 	Rating: 0.0,
	// }

	// av, err := attributevalue.MarshalMap(item)
	// if err != nil {
	// 	log.Fatalf("Got error marshalling map: %s", err)
	// }

	// // Create item in table Movies
	// input := &dynamodb.PutItemInput{
	// 	Item:      av,
	// 	TableName: aws.String(tableName),
	// }

	// _, err = c.PutItem(context.TODO(), input)
	// if err != nil {
	// 	fmt.Println("Got error calling PutItem: %s", err)
	// }

	// year := strconv.Itoa(item.Year)

	// fmt.Println("Successfully added '" + item.Title + "' (" + year + ") to table " + tableName)
	// // snippet-end:[dynamodb.go.load_items.call]

	// exists := true
	// _, err = c.DescribeTable(
	// 	context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
	// )
	// fmt.Println(err)
	// if err != nil {
	// 	var notFoundEx *types.ResourceNotFoundException
	// 	if errors.As(err, &notFoundEx) {
	// 		fmt.Println("Table %v does not exist.\n", tableName)
	// 		err = nil
	// 	} else {
	// 		fmt.Println("Couldn't determine existence of table %v. Here's why: %v\n", tableName, err)
	// 	}
	// 	exists = false
	// }
	// fmt.Println(exists)
	payload := getCiBuildPayload(c)
	UploadPlanningData("872a7985dd8a58328dea96015b738c317039fb5a", payload)
}

// snippet-end:[dynamodb.go.load_items]

func getCiBuildPayload(client *dynamodb.Client) []CiBuildPayload {
	var payload []CiBuildPayload
	originAttr, _ := attributevalue.Marshal("Tekton")
	keyExpr := expression.Key("origin").Equal(expression.Value(originAttr))
	expr, err := expression.NewBuilder().WithKeyCondition(keyExpr).Build()
	if err != nil {
		log.Fatal(err)
	}
	query, err := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String("TektonCI"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		log.Fatal(err)
	}
	// unmarshal list of items
	err = attributevalue.UnmarshalListOfMaps(query.Items, &payload)
	if err != nil {
		log.Fatal(err)
	}
	return payload
}

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
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)
}

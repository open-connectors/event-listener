// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
// snippet-start:[dynamodb.go.load_items]
package main

// snippet-start:[dynamodb.go.load_items.imports]
import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"fmt"
	"log"
	"strconv"
)

// snippet-end:[dynamodb.go.load_items.imports]

// snippet-start:[dynamodb.go.load_items.struct]
// Create struct to hold info about new item
type Item struct {
	Year   int
	Title  string
	Plot   string
	Rating float64
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
	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	// // Create DynamoDB client
	// svc := dynamodb.New(sess)
	item := Item{
		Year:   2015,
		Title:  "The Big New Movie",
		Plot:   "Nothing happens at all.",
		Rating: 0.0,
	}

	// Add each item to Movies table:
	tableName := "Movies"

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling map: %s", err)
	}

	// Create item in table Movies
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = c.PutItem(context.TODO(), input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	year := strconv.Itoa(item.Year)

	fmt.Println("Successfully added '" + item.Title + "' (" + year + ") to table " + tableName)
	// snippet-end:[dynamodb.go.load_items.call]
}

// snippet-end:[dynamodb.go.load_items]

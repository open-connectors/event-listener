package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	_ "github.com/cloudevents/sdk-go/v2"
)

const ()

// DynoObject represents an object in dynamoDB.
// Used to represent key value data such as keys, table items...
type DynoNotation map[string]types.AttributeValue

// newclient constructs a new dynamodb client using a default configuration
// and a provided profile name (created via aws configure cmd).
func newclient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://localhost:8000"}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
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

// createTable creates a table in the client's dynamodb instance
// using the table parameters specified in input.
func createTable(c *dynamodb.Client,
	tableName string, input *dynamodb.CreateTableInput,
) error {
	var tableDesc *types.TableDescription
	table, err := c.CreateTable(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to create table `%v` with error: %v\n", tableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(c)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Failed to wait on create table `%v` with error: %v\n", tableName, err)
		}
		tableDesc = table.TableDescription
	}
	fmt.Printf("Created table `%s` with details: %v\n\n", tableName, tableDesc)

	return err
}

// listTables returns a list of table names in the client's dynamodb instance.
func listTables(c *dynamodb.Client, input *dynamodb.ListTablesInput) ([]string, error) {
	tables, err := c.ListTables(
		context.TODO(),
		&dynamodb.ListTablesInput{},
	)
	if err != nil {
		return nil, err
	}

	return tables.TableNames, nil
}

// putItem inserts an item (key + attributes) in to a dynamodb table.
func putItem(c *dynamodb.Client, tableName string, item DynoNotation) (err error) {
	_, err = c.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
	})
	if err != nil {
		return err
	}

	return nil
}

// putItems batch inserts multiple items in to a dynamodb table.
func putItems(c *dynamodb.Client, tableName string, items []DynoNotation) (err error) {
	// dynamodb batch limit is 25
	if len(items) > 25 {
		return fmt.Errorf("Max batch size is 25, attempted `%d`", len(items))
	}

	// create requests
	writeRequests := make([]types.WriteRequest, len(items))
	for i, item := range items {
		writeRequests[i] = types.WriteRequest{PutRequest: &types.PutRequest{Item: item}}
	}

	// write batch
	_, err = c.BatchWriteItem(
		context.TODO(),
		&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{tableName: writeRequests},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// getItem returns an item if found based on the key provided.
// the key could be either a primary or composite key and values map.
func getItem(c *dynamodb.Client, tableName string, key DynoNotation) (item DynoNotation, err error) {
	resp, err := c.GetItem(context.TODO(), &dynamodb.GetItemInput{Key: key, TableName: aws.String(tableName)})
	if err != nil {
		return nil, err
	}

	return resp.Item, nil //
}

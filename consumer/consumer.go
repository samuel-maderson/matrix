package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func getDataFromDynamoDB(instanceName string) error {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create a DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Define the key to fetch
	key := map[string]types.AttributeValue{
		"instance_name": &types.AttributeValueMemberS{Value: instanceName},
	}

	// Fetch the item
	output, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("matrix-table"),
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	// Check if item exists
	if output.Item == nil {
		fmt.Println("No item found")
		return nil
	}

	// Print the retrieved item
	for key, value := range output.Item {
		switch v := value.(type) {
		case *types.AttributeValueMemberS:
			fmt.Printf("%s: %s\n", key, v.Value)

		default:
			fmt.Printf("%s: unknown type\n", key)
		}
	}

	return nil
}

func main() {
	instanceName := "app_02" // Replace with the actual instance name

	err := getDataFromDynamoDB(instanceName)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

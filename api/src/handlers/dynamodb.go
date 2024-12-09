package handlers

import (
	"context"
	"fmt"
	"log"
	"matrix/api/src/models"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetDataFromDynamoDB(instanceName string) models.InstancePayload {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalln("failed to load configuration: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	key := map[string]types.AttributeValue{
		"instance_name": &types.AttributeValueMemberS{Value: instanceName},
	}

	output, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("matrix-table"),
		Key:       key,
	})
	if err != nil {
		log.Fatalln("failed to get item: %w", err)
	}

	if output.Item == nil {
		fmt.Println("No item found")
	}

	instance_num := output.Item["number_of_instances"].(*types.AttributeValueMemberS).Value
	numInstances, err := strconv.Atoi(instance_num)
	if err != nil {
		log.Printf("Error converting string to integer: %v", err)

	}

	instanceFields := models.InstancePayload{
		InstanceName:         output.Item["instance_name"].(*types.AttributeValueMemberS).Value,
		SSHKeyName:           output.Item["ssh_key_name"].(*types.AttributeValueMemberS).Value,
		NumberOfInstances:    numInstances,
		InstanceType:         output.Item["instance_type"].(*types.AttributeValueMemberS).Value,
		Launched:             output.Item["launched"].(*types.AttributeValueMemberS).Value,
		WhenInstanceStart:    output.Item["when_instance_start"].(*types.AttributeValueMemberS).Value,
		WhenInstanceShutdown: output.Item["when_instance_shutdown"].(*types.AttributeValueMemberS).Value,
	}

	return instanceFields
}

func InsertDataToDynamoDB(requestData models.RequestData) error {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	item := map[string]types.AttributeValue{
		"instance_name":          &types.AttributeValueMemberS{Value: requestData.InstanceName},
		"ssh_key_name":           &types.AttributeValueMemberS{Value: requestData.SSHKeyName},
		"number_of_instances":    &types.AttributeValueMemberS{Value: requestData.NumberOfInstances},
		"instance_type":          &types.AttributeValueMemberS{Value: requestData.InstanceType},
		"launched":               &types.AttributeValueMemberS{Value: requestData.Launched},
		"when_instance_start":    &types.AttributeValueMemberS{Value: requestData.WhenInstanceStart},
		"when_instance_shutdown": &types.AttributeValueMemberS{Value: requestData.WhenInstanceShutdown},
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("matrix-table"),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	fmt.Println("Data inserted successfully")
	return nil
}

func UpdateDataDynamoDB(InstancePayload models.InstancePayload) error {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	InstancePayload.Launched = "1"

	item := map[string]types.AttributeValue{
		"instance_name":          &types.AttributeValueMemberS{Value: InstancePayload.InstanceName},
		"ssh_key_name":           &types.AttributeValueMemberS{Value: InstancePayload.SSHKeyName},
		"number_of_instances":    &types.AttributeValueMemberN{Value: strconv.Itoa(InstancePayload.NumberOfInstances)},
		"instance_type":          &types.AttributeValueMemberS{Value: InstancePayload.InstanceType},
		"launched":               &types.AttributeValueMemberS{Value: InstancePayload.Launched},
		"when_instance_start":    &types.AttributeValueMemberS{Value: InstancePayload.WhenInstanceStart},
		"when_instance_shutdown": &types.AttributeValueMemberS{Value: InstancePayload.WhenInstanceShutdown},
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("matrix-table"),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	fmt.Println("Updated successfully")
	return nil
}

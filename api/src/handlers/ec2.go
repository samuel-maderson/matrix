package handlers

import (
	"context"
	"fmt"
	"matrix/api/src/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func LaunchEC2Instance(payload models.InstancePayload) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	client := ec2.NewFromConfig(cfg)

	tags := []types.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String(payload.InstanceName),
		},
	}

	fmt.Println("PAYLOAD", payload)
	// Launch EC2 instance
	runInstancesInput := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-005fc0f236362e99f"),      // Replace with your AMI ID
		InstanceType: types.InstanceType(payload.InstanceType), // Pass the instance type as a string
		KeyName:      aws.String(payload.SSHKeyName),
		MinCount:     aws.Int32(int32(payload.NumberOfInstances)),
		MaxCount:     aws.Int32(int32(payload.NumberOfInstances)),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags:         tags,
			},
		},
		SecurityGroupIds: []string{"sg-0f36da69270ed1d01"},       // Replace with your security group ID
		SubnetId:         aws.String("subnet-02f6d5033cc0ae5cc"), // Replace with your subnet ID
	}

	output, err := client.RunInstances(context.TODO(), runInstancesInput)
	if err != nil {
		return fmt.Errorf("failed to launch EC2 instance: %w", err)
	}

	for _, instance := range output.Instances {
		fmt.Printf("Launched EC2 Instance with ID: %s\n", *instance.InstanceId)
	}
	return nil
}

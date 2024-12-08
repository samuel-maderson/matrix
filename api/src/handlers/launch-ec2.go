package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type InstancePayload struct {
	InstanceName         string `json:"instance_name"`
	SSHKeyName           string `json:"ssh_key_name"`
	NumberOfInstances    int    `json:"number_of_instances"`
	InstanceType         string `json:"instance_type"`
	WhenInstanceStart    string `json:"when_instance_start"`
	WhenInstanceShutdown string `json:"when_instance_shutdown"`
}

func launchEC2Instance(payload InstancePayload) error {
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

	// Launch EC2 instance
	runInstancesInput := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-0c02fb55956c7d316"),      // Replace with your AMI ID
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
		SecurityGroupIds: []string{"sg-XXXXX"},       // Replace with your security group ID
		SubnetId:         aws.String("subnet-XXXXX"), // Replace with your subnet ID
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

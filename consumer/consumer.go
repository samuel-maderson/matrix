package main

import (
	"fmt"

	"matrix/api/src/handlers"
	"matrix/api/src/models"
)

var (
	err  error
	data models.InstancePayload
)

func main() {
	// manual testing
	instanceName := "APP02"
	launched := 0

	if launched == 0 {
		data = handlers.GetDataFromDynamoDB(instanceName)

		err = handlers.LaunchEC2Instance(data)
		if err != nil {
			fmt.Println("Error launching EC2 instance:", err)
			return
		}

		err = handlers.UpdateDataDynamoDB(data)
		if err != nil {
			fmt.Println("Error updating data in DynamoDB:", err)
			return
		}

		fmt.Println("Consumer launched")
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type RequestData struct {
	InstanceName         string `json:"instance_name"`
	SSHKeyName           string `json:"ssh_key_name"`
	NumberOfInstances    string `json:"number_of_instances"`
	InstanceType         string `json:"instance_type"`
	Launched             string `json:"launched"`
	WhenInstanceStart    string `json:"when_instance_start"`
	WhenInstanceShutdown string `json:"when_instance_shutdown"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

var (
	PORT = "8000"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleHome(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "<h1> Welcome to the home page!</h1>")
	fmt.Fprintln(w, "<p>/ - This route is the home route<p/>")
	fmt.Fprintln(w, "<p>/launch-instance - This route launches an EC2 instance based on formdata send<p/>")
	fmt.Fprintln(w, "<p>/prepare-instance - This route stores the data will later be used to launch an instance<p/>")
	fmt.Fprintln(w, "<p>/destroy-instance - This route destroys an EC2 instance<p/>")

}

func handleDestroyInstance(w http.ResponseWriter, r *http.Request) {

	// how to respond with a json message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := ResponseMessage{
		Message: "Instance destroy request received successfully",
	}

	json.NewEncoder(w).Encode(response)
}

func handlePrepareInstance(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData RequestData

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received data: %+v\n", requestData)

	insertDataToDynamoDB(requestData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Instance launch request received successfully",
	}
	json.NewEncoder(w).Encode(response)

}

func insertDataToDynamoDB(requestData RequestData) error {

	// Initialize AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create a DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Define the item to insert
	item := map[string]types.AttributeValue{
		"instance_name":          &types.AttributeValueMemberS{Value: requestData.InstanceName},
		"ssh_key_name":           &types.AttributeValueMemberS{Value: requestData.SSHKeyName},
		"number_of_instances":    &types.AttributeValueMemberS{Value: requestData.NumberOfInstances},
		"instance_type":          &types.AttributeValueMemberS{Value: requestData.InstanceType},
		"launched":               &types.AttributeValueMemberS{Value: requestData.Launched},
		"when_instance_start":    &types.AttributeValueMemberS{Value: requestData.WhenInstanceStart},
		"when_instance_shutdown": &types.AttributeValueMemberS{Value: requestData.WhenInstanceShutdown},
	}

	// Insert the item into the table
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

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/prepare-instance", handlePrepareInstance)
	mux.HandleFunc("/destroy-instance", handleDestroyInstance)

	fmt.Println("Server is running on port:", PORT, "...")
	log.Fatal(http.ListenAndServe(":"+PORT, corsMiddleware(mux)))
}

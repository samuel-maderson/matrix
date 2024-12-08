package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestData struct {
	InstanceName         string `json:"instance_name"`
	SSHKeyName           string `json:"ssh_key_name"`
	NumberOfInstances    string `json:"number_of_instances"`
	InstanceType         string `json:"instance_type"`
	WhenInstanceStart    string `json:"when_instance_start"`
	WhenInstanceShutdown string `json:"when_instance_shutdown"`
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

func handleLaunchInstance(w http.ResponseWriter, r *http.Request) {

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Instance launch request received successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/launch-instance", handleLaunchInstance)

	fmt.Println("Server is running on port:", PORT, "...")
	log.Fatal(http.ListenAndServe(":"+PORT, corsMiddleware(mux)))
}

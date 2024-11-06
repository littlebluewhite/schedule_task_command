package main

import (
	"context"
	"github.com/littlebluewhite/schedule_task_command/proto/grpc_task_template"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Establish a connection to the server with grpc.DialContext and grpc.WithTransportCredentials
	conn, err := grpc.NewClient("localhost:55487", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create a client for the TaskService
	client := grpc_task_template.NewTaskTemplateServiceClient(conn)

	// Call an RPC method
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a request object
	request := &grpc_task_template.SendTaskTemplateRequest{
		TemplateId: 123,
		Source:     "example_source",
		TriggerFrom: []*grpc_task_template.TriggerFrom{
			{KeyValue: map[string]string{"key1": "value1", "key2": "value2"}},
		},
		TriggerAccount: "example_account",
		Variables: map[int64]*grpc_task_template.Variables{
			1: {KeyValue: map[string]string{"var1": "value1"}},
		},
	}

	// Example: Call a gRPC method (adapt this based on your service)
	response, err := client.SendTaskTemplate(ctx, request)
	if err != nil {
		log.Fatalf("Error making gRPC call: %v", err)
	}
	log.Printf("Response: %v", response)

	log.Println("Client connected and request ready")
}

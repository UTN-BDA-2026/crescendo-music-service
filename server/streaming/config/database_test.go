package config

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestConnectDB(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database connection test in short mode")
	}

	err := LoadEnv()
	if err != nil {
		t.Logf("Warning: Could not load .env file: %v", err)
	}

	dbUser := os.Getenv("MONGODB_USER")
	dbPass := os.Getenv("MONGODB_PASSWORD")
	dbHost := os.Getenv("MONGODB_HOST")
	dbPort := os.Getenv("MONGODB_PORT")

	if dbHost == "" || dbPort == "" {
		t.Fatal("Missing required .env variables")
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)

	client, err := ConnectDB(uri)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB, error: %v", err)
	}

	if client == nil {
		t.Fatal("Expected a valid MongoDB client, got nil")
	}

	defer client.Disconnect(context.TODO())
}

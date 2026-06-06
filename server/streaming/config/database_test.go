package config

import (
	"context"
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

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		t.Fatal("MONGODB_URI environment variable is not set")
	}

	client, err := ConnectDB(uri)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB, error: %v", err)
	}

	if client == nil {
		t.Fatal("Expected a valid MongoDB client, got nil")
	}

	defer client.Disconnect(context.TODO())
}

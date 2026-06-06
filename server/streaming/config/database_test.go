package config

import (
	"context"
	"testing"
)

func TestConnectDB(t *testing.T) {
	uri := "mongodb://admin:admin123@localhost:27017"

	client, err := ConnectDB(uri)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB, error: %v", err)
	}

	if client == nil {
		t.Fatal("Expected a valid MongoDB client, got nil")
	}

	defer client.Disconnect(context.TODO())
}

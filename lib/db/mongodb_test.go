package db

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestMongoConnect(t *testing.T) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		t.Skip("MONGO_URI is not set")
	}
	m := &MongoDB{
		URI:    uri,
		DBName: "admin",
	}

	m.CreateClient()
	if m.Clinet == nil {
		t.Fatalf("mongo client is nil")
	}
	defer func() {
		_ = m.Clinet.Disconnect(context.Background())
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.Clinet.Ping(ctx, readpref.Primary()); err != nil {
		t.Fatalf("ping failed: %v", err)
	}
}

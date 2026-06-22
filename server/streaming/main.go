package main

import (
	"crescendo-streaming/config"
	"crescendo-streaming/controllers"
	"crescendo-streaming/repositories"
	"crescendo-streaming/router"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Printf("Warn: Couldn't load .env: %v", err)
	}

	// Check for MONGODB_URI first (Atlas), otherwise build from individual vars (local)
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		dbUser := os.Getenv("MONGODB_USER")
		dbPass := os.Getenv("MONGODB_PASSWORD")
		dbHost := os.Getenv("MONGODB_HOST")
		dbPort := os.Getenv("MONGODB_PORT")
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)
	}

	client, err := config.ConnectDB(uri)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB successfully")

	db := client.Database("crescendo_audio")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		log.Fatalf("Failed to create GridFS bucket: %v", err)
	}

	// Initialize Redis connection (optional)
	var cacheRepo *repositories.CacheRepository
	redisClient, err := config.ConnectRedis()
	if err != nil {
		log.Printf("Warn: Couldn't connect to Redis: %v (caching disabled)", err)
	} else {
		fmt.Println("Connected to Redis successfully")
		cacheRepo = repositories.NewCacheRepository(redisClient)
	}

	streamingCtrl := controllers.NewStreamingControllerWithCache(bucket, cacheRepo)
	r := router.SetupRouter(streamingCtrl)

	port := os.Getenv("STREAMING_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

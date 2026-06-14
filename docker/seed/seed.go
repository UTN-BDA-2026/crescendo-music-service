package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UploadedFile struct {
	SongID int
	FileID string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// =========================
	// POSTGRES
	// =========================
	db := connectPostgres()
	defer db.Close()

	// =========================
	// MONGO
	// =========================
	client := connectMongo(ctx)
	defer client.Disconnect(ctx)

	bucket, _ := gridfs.NewBucket(client.Database("crescendo_audio"))

	// =========================
	// 1. UPLOAD MOCK FILES
	// =========================
	files := uploadMockFiles(ctx, bucket)

	fmt.Println("Uploaded files:", len(files))

	// =========================
	// 2. GET ALL SONGS
	// =========================
	rows, err := db.Query(`SELECT id FROM songs ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var songIDs []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		songIDs = append(songIDs, id)
	}

	// =========================
	// 3. ASSIGN FILES ROUND-ROBIN
	// =========================
	i := 0
	for _, songID := range songIDs {

		file := files[i%len(files)]

		_, err := db.Exec(`
			UPDATE songs
			SET file_id = $1
			WHERE id = $2
		`, file, songID)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Song %d -> File %s\n", songID, file)

		i++
	}

	fmt.Println("Seed completed")
}

// =========================
// UPLOAD FILES TO MONGO
// =========================
func uploadMockFiles(_ context.Context, bucket *gridfs.Bucket) []string {
	dir := "/audio"

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var fileIDs []string

	for _, e := range entries {

		path := filepath.Join(dir, e.Name())

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		fileID := primitive.NewObjectID()

		uploadStream, err := bucket.OpenUploadStreamWithID(
			fileID,
			e.Name(),
		)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(uploadStream, file)
		if err != nil {
			log.Fatal(err)
		}

		uploadStream.Close()
		file.Close()

		fileIDs = append(fileIDs, fileID.Hex())
		fmt.Println("Uploaded:", e.Name(), fileID.Hex())
	}

	return fileIDs
}

// =========================
// POSTGRES CONNECT
// =========================
func connectPostgres() *sql.DB {
	pgHost := getenv("POSTGRES_HOST", "postgresql")
	pgUser := getenv("POSTGRES_USER", "postgres")
	pgPass := getenv("POSTGRES_PASSWORD", "postgres")
	pgDB := getenv("POSTGRES_DB", "music")

	conn := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgUser, pgPass, pgDB,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// =========================
// MONGO CONNECT
// =========================
func connectMongo(ctx context.Context) *mongo.Client {
	dbUser := os.Getenv("MONGODB_USER")
	dbPass := os.Getenv("MONGODB_PASSWORD")
	dbHost := os.Getenv("MONGODB_HOST")
	dbPort := os.Getenv("MONGODB_PORT")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// =========================
// ENV
// =========================
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

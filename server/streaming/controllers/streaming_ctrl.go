package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"crescendo-streaming/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StreamingController struct {
	Bucket *gridfs.Bucket
	Cache  *repositories.CacheRepository
}

func NewStreamingController(bucket *gridfs.Bucket) *StreamingController {
	return &StreamingController{Bucket: bucket}
}

func NewStreamingControllerWithCache(bucket *gridfs.Bucket, cache *repositories.CacheRepository) *StreamingController {
	return &StreamingController{Bucket: bucket, Cache: cache}
}

func (sc *StreamingController) UploadAudio(c *gin.Context) {
	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}
	defer file.Close()

	if sc.Bucket == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Simulated file uploaded succesfully",
			"file_id": primitive.NewObjectID().Hex(),
		})
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "audio/mpeg"
	}

	opts := options.GridFSUpload().SetMetadata(bson.D{{Key: "contentType", Value: contentType}})

	uploadStream, err := sc.Bucket.OpenUploadStream(header.Filename, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload stream"})
		return
	}
	defer uploadStream.Close()

	if _, err := io.Copy(uploadStream, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		return
	}

	fileID := uploadStream.FileID.(primitive.ObjectID).Hex()

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"file_id": fileID,
	})
}

func (sc *StreamingController) StreamAudio(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured on server"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims format"})
		return
	}

	if claims["type"] != "stream" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token type"})
		return
	}

	fileIDStr, ok := claims["file_id"].(string)
	if !ok || fileIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token missing file_id"})
		return
	}

	songIDStr := ""
	if songID, ok := claims["song_id"].(string); ok {
		songIDStr = songID
	}

	if sc.Bucket == nil {
		c.Status(http.StatusOK)
		c.Writer.Write([]byte("Simulated Audio Stream"))
		return
	}

	ctx := context.Background()

	// Try to get from cache first
	if sc.Cache != nil {
		cachedStream, err := sc.Cache.GetAudioCache(ctx, fileIDStr)
		if err == nil && cachedStream != nil {
			c.Header("Content-Type", "audio/mpeg")
			c.Header("Accept-Ranges", "bytes")
			c.Status(http.StatusOK)
			io.Copy(c.Writer, cachedStream)
			cachedStream.Close()
			return
		}
	}

	// Get from MongoDB GridFS
	fileID, err := primitive.ObjectIDFromHex(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	downloadStream, err := sc.Bucket.OpenDownloadStream(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer downloadStream.Close()

	fileSize := downloadStream.GetFile().Length

	contentType := "audio/mpeg"
	metadata := downloadStream.GetFile().Metadata
	if metadata != nil {
		if val, err := metadata.LookupErr("contentType"); err == nil {
			contentType = val.StringValue()
		}
	}

	c.Header("Content-Type", contentType)
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Length", fmt.Sprintf("%d", fileSize))

	c.Status(http.StatusOK)

	// Stream to client while simultaneously caching
	if sc.Cache != nil && songIDStr != "" {
		// Read all data from stream
		data, err := io.ReadAll(downloadStream)
		if err == nil {
			// Try to cache in background (non-blocking)
			go func() {
				cacheStream := io.NopCloser(bytes.NewReader(data))
				_ = sc.Cache.SetAudioCache(context.Background(), fileIDStr, cacheStream)
			}()
			// Write to client
			c.Writer.Write(data)
		} else {
			io.Copy(c.Writer, downloadStream)
		}
	} else {
		io.Copy(c.Writer, downloadStream)
	}
}

// PauseStream handles pause requests for playback
// Expects POST request with song_id and optional position (in milliseconds)
func (sc *StreamingController) PauseStream(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured on server"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims format"})
		return
	}

	songIDStr, ok := claims["song_id"].(string)
	if !ok || songIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token missing song_id"})
		return
	}

	// Get position from request body
	var req struct {
		Position int64 `json:"position"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.Position < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Position cannot be negative"})
		return
	}

	ctx := context.Background()

	// Store pause state in cache
	if sc.Cache != nil {
		if err := sc.Cache.SetStreamingState(ctx, songIDStr, req.Position); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save playback state"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Playback paused",
		"song_id":  songIDStr,
		"position": req.Position,
	})
}

// ResumeStream handles resume requests for playback
// Retrieves the last known position for a song
func (sc *StreamingController) ResumeStream(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured on server"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims format"})
		return
	}

	songIDStr, ok := claims["song_id"].(string)
	if !ok || songIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token missing song_id"})
		return
	}

	ctx := context.Background()

	// Retrieve last known position
	var position int64 = 0
	if sc.Cache != nil {
		pos, err := sc.Cache.GetStreamingState(ctx, songIDStr)
		if err == nil {
			position = pos
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Resume information retrieved",
		"song_id":  songIDStr,
		"position": position,
	})
}

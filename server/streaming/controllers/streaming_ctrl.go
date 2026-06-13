package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StreamingController struct {
	Bucket *gridfs.Bucket
}

func NewStreamingController(bucket *gridfs.Bucket) *StreamingController {
	return &StreamingController{Bucket: bucket}
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

	if sc.Bucket == nil {
		c.Status(http.StatusOK)
		c.Writer.Write([]byte("Simulated Audio Stream"))
		return
	}

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

	io.Copy(c.Writer, downloadStream)
}

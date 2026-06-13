package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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

	uploadStream, err := sc.Bucket.OpenUploadStream(header.Filename)
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
	fileIDStr := c.Param("file_id")

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

	c.Header("Content-Type", "audio/mpeg")
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Length", fmt.Sprintf("%d", fileSize))

	c.Status(http.StatusOK)

	io.Copy(c.Writer, downloadStream)
}

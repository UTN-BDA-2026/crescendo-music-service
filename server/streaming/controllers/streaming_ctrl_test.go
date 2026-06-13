package controllers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUploadAudio(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := NewStreamingController(nil)

	r := gin.Default()
	r.POST("/upload", ctrl.UploadAudio)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audio", "test.mp3")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write([]byte("fake audio content"))
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

}

func TestStreamAudio(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := NewStreamingController(nil)

	r := gin.Default()
	r.GET("/stream/:file_id", ctrl.StreamAudio)

	req, _ := http.NewRequest(http.MethodGet, "/stream/507f1f77bcf86cd799439011", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Esperaba código %d, pero obtuve %d", http.StatusOK, w.Code)
	}
}

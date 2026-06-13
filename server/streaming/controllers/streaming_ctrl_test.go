package controllers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	t.Setenv("JWT_SECRET", "secret-key")

	ctrl := NewStreamingController(nil)

	r := gin.Default()
	r.GET("/stream", ctrl.StreamAudio)

	claims := jwt.MapClaims{
		"song_id": 1,
		"file_id": "507f1f77bcf86cd799439011",
		"type":    "stream",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "/stream?token="+tokenString, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Esperaba código %d, pero obtuve %d", http.StatusOK, w.Code)
	}
}

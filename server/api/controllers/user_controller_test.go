package controllers_test

import (
	"bytes"
	"crescendo-api/config/app"
	"crescendo-api/controllers"
	"crescendo-api/mapping"
	"crescendo-api/models"
	"crescendo-api/router"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockUserService struct{}

func (m *MockUserService) Register(req mapping.UserRegisterDTO) (models.User, error) {
	return models.User{}, nil
}

func (m *MockUserService) Login(req mapping.UserLoginDTO) (string, error) {
	return "fake-token", nil
}

func TestCanReceiveRegisterRequest(t *testing.T) {
	service := &MockUserService{}
	testRouter := router.NewRouter(&app.Container{
		User: controllers.NewUserController(service),
	})
	body := `{
				"username": "Username",
				"email":"test@mail.com",
				"date_of_birth":"1999-05-04T00:00:00Z",
				"password":"123456dvf"
			}`

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestCanReceiveLoginRequest(t *testing.T) {
	service := &MockUserService{}
	controller := controllers.NewUserController(service)

	router := gin.Default()
	router.POST("/users/login", controller.Login)

	body := `{
		"username": "Username",
		"password": "123456dvf"
	}`

	req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var response map[string]string

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	token := response["token"]

	if token != "fake-token" {
		t.Errorf("expected fake-token, got %s", token)
	}
}

package controllers_test

import (
	"bytes"
	"crescendo-api/controllers"
	"crescendo-api/mapping"
	"crescendo-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockUserService struct{}

func (m *MockUserService) Register(req mapping.UserRegisterDTO) (models.User, error) {
	return models.User{}, nil
}

func TestCanReceiveRegisterRequest(t *testing.T) {
	service := &MockUserService{}
	controller := controllers.NewUserController(service)
	router := gin.Default()
	router.POST("/users", controller.Register)
	body := `{
				"username": "Username",
				"email":"test@mail.com",
				"date_of_birth":"1999-05-04T00:00:00Z",
				"password":"123456dvf"
			}`

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

package controllers

import (
	"crescendo-api/mapping"
	"crescendo-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

func (uc *UserController) Register(c *gin.Context) {
	req := mapping.UserRegisterDTO{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	_, err := uc.service.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func (uc *UserController) Login(c *gin.Context) {
	req := mapping.UserLoginDTO{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	token, err := uc.service.Login(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

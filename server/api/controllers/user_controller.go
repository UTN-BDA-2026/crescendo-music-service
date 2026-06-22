package controllers

import (
	"crescendo-api/mapping"
	"crescendo-api/services"
	"log"
	"net/http"
	"strings"

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
		log.Printf("Register: JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	_, err := uc.service.Register(req)
	if err != nil {
		log.Printf("Register: service error: %v", err)

		errMsg := err.Error()

		// Duplicate key constraint from PostgreSQL
		if strings.Contains(errMsg, "duplicate key") || strings.Contains(errMsg, "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "El usuario o email ya existe"})
			return
		}

		// Validation errors from the service
		if strings.HasPrefix(errMsg, "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func (uc *UserController) Login(c *gin.Context) {
	req := mapping.UserLoginDTO{}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Login: JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	token, err := uc.service.Login(req)
	if err != nil {
		log.Printf("Login: service error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

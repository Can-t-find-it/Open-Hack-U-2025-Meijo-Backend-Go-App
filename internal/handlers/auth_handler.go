package handlers

import (
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service: service.AuthService{},
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dtos.LoginInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}

	response, err := h.service.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

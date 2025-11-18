package router

import (
    "github.com/gin-gonic/gin"
    "hacku_2025_meijo/backend/internal/handlers"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/health", handlers.HealthCheck)

    return r
}

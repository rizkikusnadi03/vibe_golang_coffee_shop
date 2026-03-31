package handler

import (
	"belajar_golang/pkg/response"

	"github.com/gin-gonic/gin"
)

// NewRouter membuat dan mengkonfigurasi Gin router.
func NewRouter(appEnv string) *gin.Engine {
	// Set gin mode berdasarkan environment
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Middleware bawaan
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Route group /api/v1
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			response.OK(c, "server is running", nil)
		})
	}

	return router
}

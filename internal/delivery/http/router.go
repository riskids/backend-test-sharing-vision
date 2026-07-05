package http

import (
	"article-microservice/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CORS middleware to allow frontend requests
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// SetupRouter configures all HTTP routes for the application
func SetupRouter(svc service.ArticleService, validate *validator.Validate) *gin.Engine {
	router := gin.Default()
	router.Use(corsMiddleware())

	// Create handler
	handler := NewArticleHandler(svc, validate)

	// Article routes - group under /article
	articleGroup := router.Group("/article")
	{
		// POST /article/ - Create a new article
		articleGroup.POST("/", handler.CreateHandler)

		// GET /article/:limit/:offset - Get all articles with pagination
		articleGroup.GET("/:limit/:offset", handler.GetAllHandler)

		// GET /article/detail/:id - Get article by ID
		articleGroup.GET("/detail/:id", handler.GetByIDHandler)

		// PUT /article/:id - Update article by ID
		articleGroup.PUT("/:id", handler.UpdateHandler)

		// DELETE /article/:id - Delete article by ID
		articleGroup.DELETE("/:id", handler.DeleteHandler)
	}

	return router
}
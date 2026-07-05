package http

import (
	"errors"
	"net/http"
	"strconv"

	"article-microservice/internal/delivery/http/dto"
	"article-microservice/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ArticleHandler handles HTTP requests for articles
type ArticleHandler struct {
	service   service.ArticleService
	validator *validator.Validate
}

// NewArticleHandler creates a new instance of ArticleHandler
func NewArticleHandler(service service.ArticleService, validate *validator.Validate) *ArticleHandler {
	return &ArticleHandler{
		service:   service,
		validator: validate,
	}
}

// CreateHandler handles POST /article/
func (h *ArticleHandler) CreateHandler(c *gin.Context) {
	var req dto.CreateArticleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request with custom error messages for validation requirement
	if err := h.validator.Struct(&req); err != nil {
		// Return specific validation error messages
		errors := h.formatValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	article, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return HTTP 200/201 with empty JSON body {} as per requirement
	c.JSON(http.StatusOK, gin.H{
		"id":       article.ID,
		"title":    article.Title,
		"content":  article.Content,
		"category": article.Category,
		"status":   article.Status,
	})
}

// GetAllHandler handles GET /article/:limit/:offset
func (h *ArticleHandler) GetAllHandler(c *gin.Context) {
	limitStr := c.Param("limit")
	offsetStr := c.Param("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	articles, err := h.service.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	var responses []dto.ArticleResponse
	for _, article := range articles {
		responses = append(responses, dto.ToArticleResponse(
			article.ID,
			article.Title,
			article.Content,
			article.Category,
			article.CreatedDate,
			article.UpdatedDate,
			article.Status,
		))
	}

	if responses == nil {
		responses = []dto.ArticleResponse{}
	}

	c.JSON(http.StatusOK, responses)
}

// GetByIDHandler handles GET /article/:id
func (h *ArticleHandler) GetByIDHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	article, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("article not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToArticleResponse(
		article.ID,
		article.Title,
		article.Content,
		article.Category,
		article.CreatedDate,
		article.UpdatedDate,
		article.Status,
	)

	c.JSON(http.StatusOK, response)
}

// UpdateHandler handles PUT /article/:id
func (h *ArticleHandler) UpdateHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	var req dto.UpdateArticleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		errors := h.formatValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return HTTP 200 with empty JSON body {} as per requirement
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteHandler handles DELETE /article/:id
func (h *ArticleHandler) DeleteHandler(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return HTTP 200 with empty JSON body {} as per requirement
	c.JSON(http.StatusOK, gin.H{})
}

// formatValidationErrors formats validator errors for readable output
func (h *ArticleHandler) formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := fieldErr.Field()
			tag := fieldErr.Tag()

			switch tag {
			case "required":
				errors[field] = field + " is required"
			case "min":
				if field == "Title" {
					errors[field] = "Title must be at least 20 characters"
				} else if field == "Content" {
					errors[field] = "Content must be at least 200 characters"
				} else if field == "Category" {
					errors[field] = "Category must be at least 3 characters"
				} else {
					errors[field] = field + " must be at least " + fieldErr.Param() + " characters"
				}
			case "oneof":
				errors[field] = field + " must be one of: publish, draft, thrash"
			default:
				errors[field] = field + " validation failed on " + tag
			}
		}
	}

	return errors
}
package service

import (
	"context"

	"article-microservice/internal/delivery/http/dto"
	"article-microservice/internal/model"
	"article-microservice/internal/repository"
)

// ArticleService defines the interface for article business logic
type ArticleService interface {
	Create(ctx context.Context, req dto.CreateArticleRequest) (*model.Article, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Article, error)
	GetByID(ctx context.Context, id int) (*model.Article, error)
	Update(ctx context.Context, id int, req dto.UpdateArticleRequest) error
	Delete(ctx context.Context, id int) error
}

// articleServiceImpl implements ArticleService
type articleServiceImpl struct {
	repo repository.ArticleRepository
}

// NewArticleService creates a new instance of ArticleService
func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleServiceImpl{repo: repo}
}

// Create creates a new article
func (s *articleServiceImpl) Create(ctx context.Context, req dto.CreateArticleRequest) (*model.Article, error) {
	article := &model.Article{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}

	if err := s.repo.Save(ctx, article); err != nil {
		return nil, err
	}

	return article, nil
}

// GetAll retrieves all articles with pagination
func (s *articleServiceImpl) GetAll(ctx context.Context, limit, offset int) ([]model.Article, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

// GetByID retrieves an article by ID
func (s *articleServiceImpl) GetByID(ctx context.Context, id int) (*model.Article, error) {
	return s.repo.FindById(ctx, id)
}

// Update updates an existing article
func (s *articleServiceImpl) Update(ctx context.Context, id int, req dto.UpdateArticleRequest) error {
	article := &model.Article{
		ID:       id,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}

	return s.repo.Update(ctx, article)
}

// Delete removes an article by ID
func (s *articleServiceImpl) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
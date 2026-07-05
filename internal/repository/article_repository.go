package repository

import (
	"context"
	"database/sql"
	"errors"

	"article-microservice/internal/model"
)

// ArticleRepository defines the interface for article data access
type ArticleRepository interface {
	Save(ctx context.Context, article *model.Article) error
	FindAll(ctx context.Context, limit, offset int) ([]model.Article, error)
	FindById(ctx context.Context, id int) (*model.Article, error)
	Update(ctx context.Context, article *model.Article) error
	Delete(ctx context.Context, id int) error
}

// articleRepositoryImpl implements ArticleRepository
type articleRepositoryImpl struct {
	db *sql.DB
}

// NewArticleRepository creates a new instance of ArticleRepository
func NewArticleRepository(db *sql.DB) ArticleRepository {
	return &articleRepositoryImpl{db: db}
}

// Save inserts a new article into the database
func (r *articleRepositoryImpl) Save(ctx context.Context, article *model.Article) error {
	query := `INSERT INTO posts (title, content, category, status) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, article.Title, article.Content, article.Category, article.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	article.ID = int(id)
	return nil
}

// FindAll retrieves all articles with pagination
func (r *articleRepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]model.Article, error) {
	query := `SELECT id, title, content, category, created_date, updated_date, status FROM posts LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.Category, &article.CreatedDate, &article.UpdatedDate, &article.Status)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	if articles == nil {
		articles = []model.Article{}
	}

	return articles, nil
}

// FindById retrieves an article by its ID
func (r *articleRepositoryImpl) FindById(ctx context.Context, id int) (*model.Article, error) {
	query := `SELECT id, title, content, category, created_date, updated_date, status FROM posts WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var article model.Article
	err := row.Scan(&article.ID, &article.Title, &article.Content, &article.Category, &article.CreatedDate, &article.UpdatedDate, &article.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("article not found")
		}
		return nil, err
	}

	return &article, nil
}

// Update updates an existing article in the database
func (r *articleRepositoryImpl) Update(ctx context.Context, article *model.Article) error {
	query := `UPDATE posts SET title = ?, content = ?, category = ?, status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, article.Title, article.Content, article.Category, article.Status, article.ID)
	return err
}

// Delete removes an article from the database
func (r *articleRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
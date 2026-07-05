package dto

import "time"

// ArticleResponse represents the standardized JSON response for an article
type ArticleResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Category    string    `json:"category"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
	Status      string    `json:"status"`
}

// ToArticleResponse converts model.Article to dto.ArticleResponse
func ToArticleResponse(id int, title, content, category string, createdDate, updatedDate time.Time, status string) ArticleResponse {
	return ArticleResponse{
		ID:          id,
		Title:       title,
		Content:     content,
		Category:    category,
		CreatedDate: createdDate,
		UpdatedDate: updatedDate,
		Status:      status,
	}
}

// ToArticleResponses converts slice of models to slice of dto.ArticleResponse
func ToArticleResponses(articles []ArticleResponse) []ArticleResponse {
	return articles
}
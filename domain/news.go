package domain

import (
	"time"
)

type NewsStore interface {
	Create(m NewsInput) (*News, error)
	GetManyPaginated(pp *ParsedPaginationParams) ([]*News, *Pagination, error)
	Update(m NewsInputUpdate, id int) (*News, error)
}

type News struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Categories *[]int    `json:"categories"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type NewsDB struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	UpdatedAt time.Time `db:"updated_at"`
}

type NewsInput struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Categories *[]int `json:"categories"`
}

type NewsInputUpdate struct {
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	Categories *[]int  `json:"categories"`
}

func NewsDBtoNews(newsDB *NewsDB, categories []int) *News {
	return &News{
		ID:         newsDB.ID,
		Title:      newsDB.Title,
		Content:    newsDB.Content,
		Categories: &categories,
		UpdatedAt:  newsDB.UpdatedAt,
	}
}

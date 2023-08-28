package news

import (
	"context"
)

type News struct {
	Id          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Content     string `db:"content" json:"content"`
	Description string `db:"description" json:"description"`
	Image       string `db:"image" json:"image"`
	Category    string `db:"category" json:"category"`
	CreatedAt   string `db:"created_at" json:"created_at"`
}

type CreateNewsRequest struct {
	Title       string `db:"title" json:"title"`
	Content     string `db:"content" json:"content"`
	Description string `db:"description" json:"description"`
	Image       string `db:"image" json:"image"`
	Category    string `db:"category" json:"category"`
}

type UpdateNewsRequest struct {
	Id          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Content     string `db:"content" json:"content"`
	Description string `db:"description" json:"description"`
	Image       string `db:"image" json:"image"`
	Category    string `db:"category" json:"category"`
	CreatedAt   string `db:"created_at" json:"created_at"`
}

type NewsResponse struct {
	Id          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Content     string `db:"content" json:"content"`
	Description string `db:"description" json:"description"`
	Image       string `db:"image" json:"image"`
	Category    string `db:"category" json:"category"`
	CreatedAt   string `db:"created_at" json:"created_at"`
}

type Repository interface {
	GetNews(ctx context.Context, types string) (*[]News, error)
	GetNewsByID(ctx context.Context, id int) (*News, error)
	CreateNews(ctx context.Context, news *News) error
	UpdateNews(ctx context.Context, news *News, id int) error
	DeleteNews(ctx context.Context, id int) error
	GetAllNewsCount(ctx context.Context) (int, error)
}

type Service interface {
	ViewNews(ctx context.Context, types string) (*[]NewsResponse, error)
	ViewNewsByID(ctx context.Context, id int) (*NewsResponse, error)
	AddNews(ctx context.Context, req *CreateNewsRequest) error
	ChangeNews(ctx context.Context, req *UpdateNewsRequest, id int) error
	RemoveNews(ctx context.Context, id int) error
	ViewAllNewsCount(ctx context.Context) (int, error)
}

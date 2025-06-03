package test

import (
	"BanjirEWS/news"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewAllNews(t *testing.T) {
	nr, err := news.Handler.ViewNews(*news.NewHandler(news.NewService(news.NewRepository(db.GetDB()))), context.Background(), "trending")
	assert.Nil(t, err)
	assert.Equal(t, 6, len(*nr))
}

func TestViewAllNewsCount(t *testing.T) {
	nr, err := news.Handler.ViewAllNewsCount(*news.NewHandler(news.NewService(news.NewRepository(db.GetDB()))), context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 13, nr)
}

func TestDisplayNewsById(t *testing.T) {
	nr, err := news.Handler.ViewNewsByID(*news.NewHandler(news.NewService(news.NewRepository(db.GetDB()))), context.Background(), 31)
	assert.Nil(t, err)
	assert.Equal(t, 31, nr.Id)
}

func TestRemoveNews(t *testing.T) {
	err := news.Repository.DeleteNews(news.NewRepository(db.GetDB()), context.Background(), 31)
	assert.Nil(t, err)
}

func TestCreateNews(t *testing.T) {
	req := news.News{
		Title:       "Berita 1",
		Content:     "Testing berita 1",
		Description: "Ini berita 1",
		Image:       "berita1.jpg",
		Category:    "main",
	}

	err := news.Repository.CreateNews(news.NewRepository(db.GetDB()), context.Background(), &req)
	assert.Nil(t, err)
}

func TestUpdateNews(t *testing.T) {
	req := news.News{
		Title:       "Berita 111",
		Content:     "Testing berita 1",
		Description: "Ini berita 1",
		Image:       "berita1.jpg",
		Category:    "main",
	}

	err := news.Repository.UpdateNews(news.NewRepository(db.GetDB()), context.Background(), &req, 45)
	assert.Nil(t, err)
}

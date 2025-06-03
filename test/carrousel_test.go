package test

import (
	"BanjirEWS/carrousel"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCarrousel(t *testing.T) {
	req := carrousel.Carrousel{
		Title: "Test Carrousel",
		Image: "carrousel.jpg",
		Desc:  "Ini carrousel 2",
	}

	err := carrousel.Repository.CreateCarrousel(carrousel.NewRepository(db.GetDB()), context.Background(), &req)
	assert.Nil(t, err)
}

func TestUpdateCarrousel(t *testing.T) {
	req := carrousel.Carrousel{
		Title: "Test Carrousel",
		Image: "carrousel.jpg",
		Desc:  "Ini carrousel 3",
	}

	err := carrousel.Repository.UpdateCarrousel(carrousel.NewRepository(db.GetDB()), context.Background(), &req, 1)
	assert.Nil(t, err)
}

func TestRemoveCarrousel(t *testing.T) {
	err := carrousel.Repository.DeleteCarrousel(carrousel.NewRepository(db.GetDB()), context.Background(), 34)
	assert.Nil(t, err)
}

func TestGetCarrouselByID(t *testing.T) {
	c, err := carrousel.Repository.GetCarrouselByID(carrousel.NewRepository(db.GetDB()), context.Background(), 30)
	assert.Nil(t, err)
	assert.Equal(t, 30, c.Id)
}

func TestGetCarrouselByIDAdmin(t *testing.T) {
	c, err := carrousel.Repository.GetCarrouselByIDAdmin(carrousel.NewRepository(db.GetDB()), context.Background(), 30)
	assert.Nil(t, err)
	assert.Equal(t, 30, c.Id)
}

func TestGetAllCarrouselCount(t *testing.T) {
	count, err := carrousel.Repository.GetAllCarrouselCount(carrousel.NewRepository(db.GetDB()), context.Background())
	assert.Nil(t, err)
	assert.Greater(t, count, 0)
}

func TestGetCarrousels(t *testing.T) {
	cs, err := carrousel.Repository.GetCarrousel(carrousel.NewRepository(db.GetDB()), context.Background())
	assert.Nil(t, err)
	assert.NotEmpty(t, cs)
}

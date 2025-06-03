package test

import (
	"BanjirEWS/gallery"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGallery(t *testing.T) {
	req := gallery.Gallery{
		Title: "Gallery 1",
		Image: "gallery1.jpg",
		Date:  "10 Januari 2023",
	}

	err := gallery.Repository.CreateGallery(gallery.NewRepository(db.GetDB()), context.Background(), &req)
	assert.Nil(t, err)
}

func TestUpdateGallery(t *testing.T) {
	req := gallery.Gallery{
		Title: "Gallery 111",
		Image: "gallery1.jpg",
		Date:  "10 Januari 2023",
	}

	err := gallery.Repository.UpdateGallery(gallery.NewRepository(db.GetDB()), context.Background(), &req, 20)
	assert.Nil(t, err)
}

func TestRemoveGallery(t *testing.T) {
	err := gallery.Repository.DeleteGallery(gallery.NewRepository(db.GetDB()), context.Background(), 1)
	assert.Nil(t, err)
}

func TestGetGalleryByID(t *testing.T) {
	gr, err := gallery.Repository.GetGalleryById(gallery.NewRepository(db.GetDB()), context.Background(), 20)
	assert.Nil(t, err)
	assert.Equal(t, 20, gr.Id)
}

func TestGetGalleryByIDAdmin(t *testing.T) {
	gr, err := gallery.Repository.GetGalleryByIdAdmin(gallery.NewRepository(db.GetDB()), context.Background(), 20)
	assert.Nil(t, err)
	assert.Equal(t, 20, gr.Id)
}

func TestGetAllGalleryCount(t *testing.T) {
	grCount, err := gallery.Repository.GetAllGalleryCount(gallery.NewRepository(db.GetDB()), context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 10, grCount)
}

func TestGetAllGallery(t *testing.T) {
	gr, err := gallery.Repository.GetGallery(gallery.NewRepository(db.GetDB()), context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 10, len(*gr))
}

package gallery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type service struct {
	Repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		Repository: repository,
	}
}

func (s *service) AddGallery(ctx context.Context, req *CreateGalleryRequest) error {
	gallery := Gallery{
		Title: req.Title,
		Image: req.Image,
		Date:  req.Date,
	}
	err := s.Repository.CreateGallery(ctx, &gallery)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ViewGallery(ctx context.Context) (*[]GalleryResponse, error) {
	var galleries []GalleryResponse
	res, err := s.Repository.GetGallery(ctx)
	if err != nil {
		return &[]GalleryResponse{}, err
	}
	for _, g := range *res {
		gallery := NewGallery(g)

		galleries = append(galleries, *gallery)
	}
	return &galleries, nil
}

func NewGallery(gallery Gallery) *GalleryResponse {
	return &GalleryResponse{
		Id:    gallery.Id,
		Title: gallery.Title,
		Image: gallery.Image,
		Date:  gallery.Date,
	}
}

func (s *service) ViewGalleryById(ctx context.Context, id int) (*GalleryResponse, error) {
	res, err := s.Repository.GetGalleryById(ctx, id)
	if err != nil {
		return &GalleryResponse{}, err
	}

	gallery := GalleryResponse{
		Id:    res.Id,
		Title: res.Title,
		Image: res.Image,
		Date:  res.Date,
	}

	return &gallery, nil

}

func (s *service) ViewGalleryByIdAdmin(ctx context.Context, id int) (*GalleryResponse, error) {
	res, err := s.Repository.GetGalleryByIdAdmin(ctx, id)
	if err != nil {
		return &GalleryResponse{}, err
	}

	gallery := GalleryResponse{
		Id:    res.Id,
		Title: res.Title,
		Image: res.Image,
		Date:  res.Date,
	}

	return &gallery, nil

}

func (s *service) RemoveGallery(ctx context.Context, id int) error {
	err := s.Repository.DeleteGallery(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("gallery with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ChangeGallery(ctx context.Context, req *UpdateGalleryRequest, id int) error {
	gallery := Gallery{
		Id:    req.Id,
		Title: req.Title,
		Image: req.Image,
		Date:  req.Date,
	}

	err := s.Repository.UpdateGallery(ctx, &gallery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("gallery with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ViewAllGalleryCount(ctx context.Context) (int, error) {
	count, err := s.Repository.GetAllGalleryCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, err
}

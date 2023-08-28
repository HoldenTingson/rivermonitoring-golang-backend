package carrousel

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type service struct {
	Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (s *service) ViewCarrousel(ctx context.Context) (*[]CarrouselResponse, error) {
	n, err := s.Repository.GetCarrousel(ctx)
	if err != nil {
		return &[]CarrouselResponse{}, err
	}

	var Carrousel []CarrouselResponse

	for _, res := range *n {
		CarrouselRes := &CarrouselResponse{
			Id:    res.Id,
			Title: res.Title,
			Desc:  res.Desc,
			Image: res.Image,
			Date:  res.Date,
		}

		Carrousel = append(Carrousel, *CarrouselRes)
	}

	return &Carrousel, nil
}

func (s *service) ViewCarrouselByID(ctx context.Context, id int) (*CarrouselResponse, error) {
	res, err := s.Repository.GetCarrouselByID(ctx, id)
	if err != nil {
		return &CarrouselResponse{}, err
	}

	Carrousel := CarrouselResponse{
		Id:    res.Id,
		Title: res.Title,
		Desc:  res.Desc,
		Image: res.Image,
		Date:  res.Date,
	}

	return &Carrousel, nil
}

func (s *service) ViewCarrouselByIDAdmin(ctx context.Context, id int) (*CarrouselResponse, error) {
	res, err := s.Repository.GetCarrouselByIDAdmin(ctx, id)
	if err != nil {
		return &CarrouselResponse{}, err
	}

	Carrousel := CarrouselResponse{
		Id:    res.Id,
		Title: res.Title,
		Desc:  res.Desc,
		Image: res.Image,
		Date:  res.Date,
	}

	return &Carrousel, nil
}

func (s *service) AddCarrousel(ctx context.Context, req *CreateCarrouselRequest) error {
	Carrousel := &Carrousel{
		Title: req.Title,
		Desc:  req.Desc,
		Image: req.Image,
		Date:  req.Date,
	}

	err := s.Repository.CreateCarrousel(ctx, Carrousel)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ChangeCarrousel(ctx context.Context, req *UpdateCarrouselRequest, id int) error {
	Carrousel := Carrousel{
		Title: req.Title,
		Desc:  req.Desc,
		Image: req.Image,
		Date:  req.Date,
	}

	err := s.Repository.UpdateCarrousel(ctx, &Carrousel, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Carrousel with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) RemoveCarrousel(ctx context.Context, id int) error {
	err := s.Repository.DeleteCarrousel(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Carrousel with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ViewAllCarrouselCount(ctx context.Context) (int, error) {
	count, err := s.Repository.GetAllCarrouselCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, err
}

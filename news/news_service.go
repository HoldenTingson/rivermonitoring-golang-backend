package news

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

func (s *service) ViewNews(ctx context.Context, types string) (*[]NewsResponse, error) {
	n, err := s.Repository.GetNews(ctx, types)
	if err != nil {
		return &[]NewsResponse{}, err
	}

	var news []NewsResponse

	for _, res := range *n {
		newsRes := &NewsResponse{
			Id:          res.Id,
			Title:       res.Title,
			Content:     res.Content,
			Description: res.Description,
			Image:       res.Image,
			Category:    res.Category,
			CreatedAt:   res.CreatedAt,
		}

		news = append(news, *newsRes)
	}

	return &news, nil
}

func (s *service) ViewNewsByID(ctx context.Context, id int) (*NewsResponse, error) {
	res, err := s.Repository.GetNewsByID(ctx, id)
	if err != nil {
		return &NewsResponse{}, err
	}

	news := NewsResponse{
		Id:          res.Id,
		Title:       res.Title,
		Content:     res.Content,
		Description: res.Description,
		Image:       res.Image,
		Category:    res.Category,
		CreatedAt:   res.CreatedAt,
	}

	return &news, nil
}

func (s *service) AddNews(ctx context.Context, req *CreateNewsRequest) error {
	news := &News{
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Image:       req.Image,
		Category:    req.Category,
	}

	err := s.Repository.CreateNews(ctx, news)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ChangeNews(ctx context.Context, req *UpdateNewsRequest, id int) error {
	news := News{

		Id:          req.Id,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Image:       req.Image,
		Category:    req.Category,
	}

	err := s.Repository.UpdateNews(ctx, &news, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("news with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) RemoveNews(ctx context.Context, id int) error {
	err := s.Repository.DeleteNews(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("news with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ViewAllNewsCount(ctx context.Context) (int, error) {
	count, err := s.Repository.GetAllNewsCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, err
}

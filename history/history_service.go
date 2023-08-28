package history

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type service struct {
	Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (s *service) ViewHistoryByRiverIdByTime(ctx context.Context, id string) (*[]HistoryResponse, error) {
	var histories []HistoryResponse
	res, err := s.Repository.GetHistoryByRiverIdByTime(ctx, id)
	if err != nil {
		return &[]HistoryResponse{}, err
	}
	for _, h := range *res {
		history := NewHistory(h)

		histories = append(histories, *history)
	}
	return &histories, nil
}

func NewHistory(history History) *HistoryResponse {
	return &HistoryResponse{
		Id:        history.Id,
		Height:    history.Height,
		Status:    history.Status,
		Timestamp: history.Timestamp,
		RiverId:   history.RiverId,
	}
}

func (s *service) ViewHistoryByRiverId(ctx context.Context, id string) (*[]HistoryResponse, error) {
	var histories []HistoryResponse
	res, err := s.Repository.GetHistoryByRiverId(ctx, id)
	if err != nil {
		return &[]HistoryResponse{}, err
	}
	for _, h := range *res {
		history := NewHistory(h)

		histories = append(histories, *history)
	}
	return &histories, nil
}

func (s *service) ViewHistoryById(ctx context.Context, id int) (*HistoryResponse, error) {
	res, err := s.Repository.GetHistoryById(ctx, id)
	if err != nil {
		return &HistoryResponse{}, err
	}

	history := NewHistory(*res)

	return history, nil
}

func (s *service) ViewHistoryCountByRiverId(ctx context.Context, id string) (int, error) {
	count, err := s.Repository.GetHistoryCountByRiverId(ctx, id)
	if err != nil {
		return 0, err
	}

	return count, err
}

func (s *service) RemoveAllHistoryByTime(ctx context.Context, startTime, endTime time.Time) error {
	err := s.Repository.DeleteAllHistoryByTime(ctx, startTime, endTime)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) RemoveAllHistory(ctx context.Context) error {
	err := s.Repository.DeleteAllHistory(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) RemoveHistoryByRiverId(ctx context.Context, id string) error {
	err := s.Repository.DeleteHistoryByRiverId(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("history with ID %s not found", id)
		}
		return err
	}

	return nil
}

func (s *service) RemoveHistoryById(ctx context.Context, id int) error {
	err := s.Repository.DeleteHistoryById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("history with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) RemoveHistoryByRiverIdByTime(ctx context.Context, id string, startTime, endTime time.Time) error {
	err := s.Repository.DeleteHistoryByRiverIdByTime(ctx, id, startTime, endTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("history with ID %s not found", id)
		}
		return err
	}

	return nil
}

package feedback

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

func (s *service) AddFeedback(ctx context.Context, req *CreateFeedbackRequest) error {
	feedback := Feedback{
		Content:  req.Content,
		ReportId: req.ReportId,
	}

	err := s.Repository.CreateFeedback(ctx, &feedback)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ViewFeedbacks(ctx context.Context) (*[]FeedbackResponse, error) {
	var feedbacks []FeedbackResponse
	res, err := s.Repository.GetFeedback(ctx)
	if err != nil {
		return &[]FeedbackResponse{}, err
	}

	feedbacks = append(feedbacks, *res...)

	return &feedbacks, nil
}

func (s *service) ViewFeedbackByReportId(ctx context.Context, id int) (*FeedbackResponse, error) {
	res, err := s.Repository.GetFeedbackByReportId(ctx, id)
	if err != nil {
		return &FeedbackResponse{}, err
	}

	return res, nil
}

func (s *service) ViewFeedbackById(ctx context.Context, id int) (*FeedbackResponse, error) {
	res, err := s.Repository.GetFeedbackById(ctx, id)
	if err != nil {
		return &FeedbackResponse{}, err
	}

	return res, nil
}

func (s *service) RemoveFeedback(ctx context.Context, id int) error {
	err := s.Repository.DeleteFeedback(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("feedback with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ChangeFeedback(ctx context.Context, req *UpdateFeedbackRequest, id int) error {
	feedback := Feedback{
		Content:  req.Content,
		ReportId: req.ReportId,
	}

	err := s.Repository.UpdateFeedback(ctx, &feedback, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("feedback with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ViewAllFeedbackCount(ctx context.Context) (int, error) {
	count, err := s.Repository.GetAllFeedbackCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, err
}

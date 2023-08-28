package report

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

func (s *service) ViewReportByUserId(ctx context.Context, id int) (*[]Report, error) {
	var reports []Report
	r, err := s.Repository.GetReporBytUserId(ctx, id)
	if err != nil {
		return &[]Report{}, err
	}

	reports = append(reports, *r...)

	return &reports, nil
}
func (s *service) ViewReportById(ctx context.Context, id int) (*ReportResponse, error) {
	res, err := s.Repository.GetReportById(ctx, id)
	if err != nil {
		return &ReportResponse{}, err
	}

	return res, nil
}

func (s *service) AddReport(ctx context.Context, req *CreateReportRequest) error {
	r := &Report{
		Title:   req.Title,
		Content: req.Content,
		UserId:  req.UserId,
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
	}

	err := s.Repository.CreateReport(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) RemoveReport(ctx context.Context, id int) error {
	err := s.Repository.DeleteReport(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("report with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ViewReports(ctx context.Context) (*[]ReportResponse, error) {
	var reports []ReportResponse
	res, err := s.Repository.GetReports(ctx)
	if err != nil {
		return &[]ReportResponse{}, err
	}
	reports = append(reports, *res...)
	return &reports, nil
}

func (s *service) ViewReportByUserIdById(ctx context.Context, id int) (*ReportResponse, error) {
	res, err := s.Repository.GetReportByUserIdById(ctx, id)
	if err != nil {
		return &ReportResponse{}, err
	}

	return res, nil
}

func (s *service) ViewAllReportCount(ctx context.Context) (int, error) {
	count, err := s.Repository.GetAllReportCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, err
}

func (s *service) ViewUserReportCount(ctx context.Context, id int) (int, error) {
	count, err := s.Repository.GetReportCountByUserId(ctx, id)
	if err != nil {
		return 0, err
	}

	return count, err
}

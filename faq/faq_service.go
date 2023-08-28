package faq

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

func (s *service) AddFaq(ctx context.Context, req *CreateFaqRequest) error {
	faq := Faq{
		Category: req.Category,
		Question: req.Question,
		Answer:   req.Answer,
	}
	err := s.Repository.CreateFaq(ctx, &faq)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) ViewFaq(ctx context.Context) (*[]FaqResponse, error) {
	var faqs []FaqResponse
	res, err := s.Repository.GetFaq(ctx)
	if err != nil {
		return &[]FaqResponse{}, err
	}
	for _, f := range *res {
		faq := NewFaq(f)

		faqs = append(faqs, *faq)
	}
	return &faqs, nil
}

func NewFaq(faq Faq) *FaqResponse {
	return &FaqResponse{
		Id:       faq.Id,
		Category: faq.Category,
		Question: faq.Question,
		Answer:   faq.Answer,
	}
}

func (s *service) RemoveFaq(ctx context.Context, id int) error {
	err := s.Repository.DeleteFaq(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("faq with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) ChangeFaq(ctx context.Context, req *UpdateFaqRequest, id int) error {
	faq := Faq{
		Category: req.Category,
		Question: req.Question,
		Answer:   req.Answer,
	}

	err := s.Repository.UpdateFaq(ctx, &faq, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("faq with ID %d not found", id)
		}
		return err
	}

	return nil
}

func (s *service) SearchFaq(ctx context.Context, question string) (*[]FaqResponse, error) {
	res, err := s.Repository.FindFaq(ctx, question)
	if err != nil {
		return nil, err
	}

	var faqs []FaqResponse
	for _, faq := range *res {
		faqRes := NewFaqResponse(faq)

		faqs = append(faqs, *faqRes)
	}

	return &faqs, nil
}

func NewFaqResponse(res Faq) *FaqResponse {
	return &FaqResponse{
		Id:       res.Id,
		Question: res.Question,
		Answer:   res.Answer,
	}
}

func (s *service) ViewCategory(ctx context.Context) (*[]FaqResponse, error) {
	var faqs []FaqResponse
	res, err := s.Repository.GetCategory(ctx)
	if err != nil {
		return &[]FaqResponse{}, err
	}
	for _, f := range *res {
		faq := NewFaq(f)

		faqs = append(faqs, *faq)
	}
	return &faqs, nil
}

func (s *service) ViewQa(ctx context.Context, category string) (*[]FaqResponse, error) {
	var faqs []FaqResponse
	res, err := s.Repository.GetQa(ctx, category)
	if err != nil {
		return &[]FaqResponse{}, err
	}
	for _, f := range *res {
		faq := NewFaq(f)

		faqs = append(faqs, *faq)
	}
	return &faqs, nil
}

func (s *service) ViewAllFaqCount(ctx context.Context) (int, error) {
	count, err := s.Repository.GetAllFaqCount(ctx)
	if err != nil {
		return 0, err
	}

	return count, err
}

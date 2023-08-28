package faq

import (
	"context"
)

type Faq struct {
	Id       int    `db:"id" json:"id"`
	Category string `db:"category" json:"category"`
	Question string `db:"question" json:"question"`
	Answer   string `db:"answer" json:"answer"`
}

type FaqResponse struct {
	Id       int    `db:"id" json:"id,omitempty"`
	Category string `db:"category" json:"category,omitempty"`
	Question string `db:"question" json:"question,omitempty"`
	Answer   string `db:"answer" json:"answer,omitempty"`
}

type UpdateFaqRequest struct {
	Category string `db:"category" json:"category"`
	Question string `db:"question" json:"question"`
	Answer   string `db:"answer" json:"answer"`
}

type CreateFaqRequest struct {
	Category string `db:"category" json:"category"`
	Question string `db:"question" json:"question"`
	Answer   string `db:"answer" json:"answer"`
}

type Repository interface {
	GetFaq(ctx context.Context) (*[]Faq, error)
	UpdateFaq(ctx context.Context, faq *Faq, id int) error
	DeleteFaq(ctx context.Context, id int) error
	CreateFaq(ctx context.Context, faq *Faq) error
	FindFaq(ctx context.Context, question string) (*[]Faq, error)
	GetCategory(ctx context.Context) (*[]Faq, error)
	GetQa(ctx context.Context, category string) (*[]Faq, error)
	GetAllFaqCount(ctx context.Context) (int, error)
}

type Service interface {
	ViewFaq(ctx context.Context) (*[]FaqResponse, error)
	ChangeFaq(ctx context.Context, req *UpdateFaqRequest, id int) error
	RemoveFaq(ctx context.Context, id int) error
	AddFaq(ctx context.Context, req *CreateFaqRequest) error
	SearchFaq(ctx context.Context, question string) (*[]FaqResponse, error)
	ViewCategory(ctx context.Context) (*[]FaqResponse, error)
	ViewQa(ctx context.Context, category string) (*[]FaqResponse, error)
	ViewAllFaqCount(ctx context.Context) (int, error)
}

package river

import (
	"context"
)

type River struct {
	Id        string  `db:"id" json:"id"`
	Latitude  string  `db:"latitude" json:"latitude"`
	Longitude string  `db:"longitude" json:"longitude"`
	Location  string  `db:"location" json:"location"`
	Height    float64 `db:"height" json:"height"`
	Status    string  `db:"status" json:"status"`
}

type UpdateRiverRequest struct {
	Id        string  `db:"id" json:"id"`
	Latitude  string  `db:"latitude" json:"latitude"`
	Longitude string  `db:"longitude" json:"longitude"`
	Location  string  `db:"location" json:"location"`
	Height    float64 `db:"height" json:"height"`
	Status    string  `db:"status" json:"status"`
}

type CreateRiverRequest struct {
	Id        string `db:"id" json:"id"`
	Latitude  string `db:"latitude" json:"latitude"`
	Longitude string `db:"longitude" json:"longitude"`
	Location  string `db:"location" json:"location"`
}

type RiverResponse struct {
	Id        string  `db:"id" json:"id"`
	Latitude  string  `db:"latitude" json:"latitude"`
	Longitude string  `db:"longitude" json:"longitude"`
	Location  string  `db:"location" json:"location"`
	Height    float64 `db:"height" json:"height"`
	Status    string  `db:"status" json:"status"`
}

type UpdateRiver struct {
	Id     string  `db:"id" json:"id"`
	Height float64 `db:"height" json:"height"`
	Status string  `db:"status" json:"status"`
}

type Repository interface {
	CreateRiver(ctx context.Context, river *River) error
	DeleteRiver(ctx context.Context, id string) error
	UpdateRiverDetail(ctx context.Context, river *River, id string) error
	GetRiver(ctx context.Context) (*[]River, error)
	GetRiverId(ctx context.Context) ([]string, error)
	GetRiverById(ctx context.Context, id string) (*River, error)
	UpdateRiver(ctx context.Context, riverChan chan *UpdateRiver, id string) error
	FindRiver(ctx context.Context, location string) (*River, error)
	FilterRiver(ctx context.Context, sortBy string) (*[]River, error)
	GetAllRiverCount(ctx context.Context) (int, error)
	GetRiverByStatus(ctx context.Context, status string) (*[]River, error)
}

type Service interface {
	AddRiver(ctx context.Context, req *CreateRiverRequest) error
	RemoveRiver(ctx context.Context, id string) error
	ChangeRiverDetail(ctx context.Context, req *UpdateRiverRequest, id string) error
	ViewRiver(ctx context.Context) (*[]RiverResponse, error)
	ViewRiverById(ctx context.Context, id string) (*RiverResponse, error)
	UpdateRiver(ctx context.Context) error
	SearchRiver(ctx context.Context, location string) (*RiverResponse, error)
	SortRiver(ctx context.Context, sortBy string) (*[]RiverResponse, error)
	ViewAllRiverCount(ctx context.Context) (int, error)
	ViewRiverByStatus(ctx context.Context, status string) (*[]RiverResponse, error)
}

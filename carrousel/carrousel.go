package carrousel

import (
	"context"
)

type Carrousel struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Desc  string `db:"description" json:"desc"`
	Image string `db:"image" json:"image"`
	Date  string `db:"date" json:"date"`
}

type CreateCarrouselRequest struct {
	Title string `db:"title" json:"title"`
	Desc  string `db:"description" json:"desc"`
	Image string `db:"image" json:"image"`
	Date  string `db:"date" json:"date"`
}

type UpdateCarrouselRequest struct {
	Title string `db:"title" json:"title"`
	Desc  string `db:"description" json:"desc"`
	Image string `db:"image" json:"image"`
	Date  string `db:"date" json:"date"`
}

type CarrouselResponse struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Desc  string `db:"description" json:"desc"`
	Image string `db:"image" json:"image"`
	Date  string `db:"date" json:"date"`
}

type Repository interface {
	GetCarrousel(ctx context.Context) (*[]Carrousel, error)
	GetCarrouselByID(ctx context.Context, id int) (*Carrousel, error)
	GetCarrouselByIDAdmin(ctx context.Context, id int) (*Carrousel, error)
	CreateCarrousel(ctx context.Context, Carrousel *Carrousel) error
	UpdateCarrousel(ctx context.Context, Carrousel *Carrousel, id int) error
	DeleteCarrousel(ctx context.Context, id int) error
	GetAllCarrouselCount(ctx context.Context) (int, error)
}

type Service interface {
	ViewCarrousel(ctx context.Context) (*[]CarrouselResponse, error)
	ViewCarrouselByID(ctx context.Context, id int) (*CarrouselResponse, error)
	ViewCarrouselByIDAdmin(ctx context.Context, id int) (*CarrouselResponse, error)
	AddCarrousel(ctx context.Context, req *CreateCarrouselRequest) error
	ChangeCarrousel(ctx context.Context, req *UpdateCarrouselRequest, id int) error
	RemoveCarrousel(ctx context.Context, id int) error
	ViewAllCarrouselCount(ctx context.Context) (int, error)
}

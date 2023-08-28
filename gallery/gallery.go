package gallery

import "context"

type Gallery struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Image string `db:"image" json:"image"`
	Date  string `db:"date" json:"date"`
}

type UpdateGalleryRequest struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Image string `db:"image" json:"image"`
	Date  string `db:"date" json:"date"`
}

type CreateGalleryRequest struct {
	Title string `db:"title" json:"title"`
	Image string `db:"image" json:"image"`
	Date  string `db:"date" json:"date"`
}

type GalleryResponse struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Image string `db:"image" json:"image"`
	Date  string `db:"date" json:"date"`
}

type Repository interface {
	GetGallery(ctx context.Context) (*[]Gallery, error)
	CreateGallery(ctx context.Context, gallery *Gallery) error
	GetGalleryById(ctx context.Context, id int) (*Gallery, error)
	GetGalleryByIdAdmin(ctx context.Context, id int) (*Gallery, error)
	DeleteGallery(ctx context.Context, id int) error
	UpdateGallery(ctx context.Context, gallery *Gallery, id int) error
	GetAllGalleryCount(ctx context.Context) (int, error)
}

type Service interface {
	ViewGallery(ctx context.Context) (*[]GalleryResponse, error)
	AddGallery(ctx context.Context, req *CreateGalleryRequest) error
	ViewGalleryById(ctx context.Context, id int) (*GalleryResponse, error)
	ViewGalleryByIdAdmin(ctx context.Context, id int) (*GalleryResponse, error)
	RemoveGallery(ctx context.Context, id int) error
	ChangeGallery(ctx context.Context, req *UpdateGalleryRequest, id int) error
	ViewAllGalleryCount(ctx context.Context) (int, error)
}

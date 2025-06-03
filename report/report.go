package report

import "context"

type Report struct {
	Id         int    `db:"id" json:"id"`
	Content    string `db:"content" json:"content"`
	Attachment string `db:"attachment" json:"attachment"`
	UserId     int    `db:"user_id" json:"user_id"`
	Name       string `db:"name" json:"name"`
	Email      string `db:"email" json:"email"`
	Phone      string `db:"phone" json:"phone"`
	CreatedAt  string `db:"created_at" json:"created_at"`
}

type ReportResponse struct {
	Id         int    `db:"id" json:"id"`
	Content    string `db:"content" json:"content"`
	Attachment string `db:"attachment" json:"attachment"`
	UserId     int    `db:"user_id" json:"user_id"`
	Name       string `db:"name" json:"name"`
	Email      string `db:"email" json:"email"`
	Phone      string `db:"phone" json:"phone"`
	CreatedAt  string `db:"created_at" json:"created_at"`
}

type CreateReportRequest struct {
	Content    string `db:"content" json:"content"`
	Attachment string `db:"attachment" json:"attachment"`
	UserId     int    `db:"user_id" json:"user_id"`
	Name       string `db:"name" json:"name"`
	Email      string `db:"email" json:"email"`
	Phone      string `db:"phone" json:"phone"`
}

type Repository interface {
	GetReporBytUserId(ctx context.Context, id int) (*[]Report, error)
	GetReportById(ctx context.Context, id int) (*ReportResponse, error)
	CreateReport(ctx context.Context, report *Report) error
	DeleteReport(ctx context.Context, id int) error
	GetReports(ctx context.Context) (*[]ReportResponse, error)
	GetReportByUserIdById(ctx context.Context, id int) (*ReportResponse, error)
	GetReportCountByUserId(ctx context.Context, id int) (int, error)
	GetAllReportCount(ctx context.Context) (int, error)
}

type Service interface {
	ViewReportByUserId(ctx context.Context, id int) (*[]Report, error)
	ViewReportById(ctx context.Context, id int) (*ReportResponse, error)
	AddReport(ctx context.Context, req *CreateReportRequest) error
	RemoveReport(ctx context.Context, id int) error
	ViewReports(ctx context.Context) (*[]ReportResponse, error)
	ViewReportByUserIdById(ctx context.Context, id int) (*ReportResponse, error)
	ViewUserReportCount(ctx context.Context, id int) (int, error)
	ViewAllReportCount(ctx context.Context) (int, error)
}

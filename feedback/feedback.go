package feedback

import "context"

type Feedback struct {
	Id       int    `db:"id" json:"id"`
	Content  string `db:"content" json:"content"`
	ReportId int    `db:"report_id" json:"report_id"`
}

type CreateFeedbackRequest struct {
	Content  string `db:"content" json:"content"`
	ReportId int    `db:"report_id" json:"report_id"`
}

type UpdateFeedbackRequest struct {
	Content  string `db:"content" json:"content"`
	ReportId int    `db:"report_id" json:"report_id"`
}

type FeedbackResponse struct {
	Id       int    `db:"id" json:"id"`
	Content  string `db:"content" json:"content"`
	Username string `db:"username" json:"username,omitempty"`
}

type Repository interface {
	CreateFeedback(ctx context.Context, feedback *Feedback) error
	GetFeedback(ctx context.Context) (*[]FeedbackResponse, error)
	GetFeedbackByReportId(ctx context.Context, id int) (*FeedbackResponse, error)
	GetFeedbackById(ctx context.Context, id int) (*FeedbackResponse, error)
	DeleteFeedback(ctx context.Context, id int) error
	UpdateFeedback(ctx context.Context, feedback *Feedback, id int) error
	GetAllFeedbackCount(ctx context.Context) (int, error)
}

type Service interface {
	AddFeedback(ctx context.Context, req *CreateFeedbackRequest) error
	ViewFeedbacks(ctx context.Context) (*[]FeedbackResponse, error)
	ViewFeedbackByReportId(ctx context.Context, id int) (*FeedbackResponse, error)
	ViewFeedbackById(ctx context.Context, id int) (*FeedbackResponse, error)
	RemoveFeedback(ctx context.Context, id int) error
	ChangeFeedback(ctx context.Context, req *UpdateFeedbackRequest, id int) error
	ViewAllFeedbackCount(ctx context.Context) (int, error)
}

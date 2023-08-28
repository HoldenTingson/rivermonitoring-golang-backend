package feedback

import (
	"context"
	"database/sql"
	"errors"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateFeedback(ctx context.Context, feedback *Feedback) error {
	query := "insert into feedback (content, report_id) values(?,?)"
	_, err := r.db.ExecContext(ctx, query, feedback.Content, feedback.ReportId)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetFeedback(ctx context.Context) (*[]FeedbackResponse, error) {
	var feedbacks []FeedbackResponse
	query := "select id, content FROM feedback"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return &[]FeedbackResponse{}, err
	}

	for rows.Next() {
		var feedback FeedbackResponse
		rows.Scan(&feedback.Id, &feedback.Content)
		feedbacks = append(feedbacks, feedback)
	}

	return &feedbacks, nil
}

func (r *repository) GetFeedbackByReportId(ctx context.Context, id int) (*FeedbackResponse, error) {
	var feedback FeedbackResponse
	query := "select id, content FROM feedback where report_id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&feedback.Id, &feedback.Content)
	if err != nil {
		return &FeedbackResponse{}, err
	}
	return &feedback, nil
}

func (r *repository) GetFeedbackById(ctx context.Context, id int) (*FeedbackResponse, error) {
	var feedback FeedbackResponse
	query := "select f.id, f.content, u.username FROM feedback f JOIN report r ON f.report_id = r.id JOIN user u ON r.user_id = u.id where f.id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&feedback.Id, &feedback.Content, &feedback.Username)
	if err != nil {
		return &FeedbackResponse{}, err
	}
	return &feedback, nil
}

func (r *repository) UpdateFeedback(ctx context.Context, feedback *Feedback, id int) error {
	query := "update feedback set content = ? where id = ?"
	res, err := r.db.ExecContext(ctx, query, feedback.Content, id)
	if err != nil {
		panic(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (r *repository) DeleteFeedback(ctx context.Context, id int) error {
	query := "delete from feedback where id = ?"
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (r *repository) GetAllFeedbackCount(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM feedback"
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

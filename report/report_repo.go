package report

import (
	"BanjirEWS/util"
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	return &repository{db: db}
}

func (r *repository) GetReporBytUserId(ctx context.Context, id int) (*[]Report, error) {
	var reports []Report
	query := "select id, content from report where user_id = ?"
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return &[]Report{}, err
	}

	for rows.Next() {
		var report Report
		err := rows.Scan(&report.Id, &report.Content)

		if err != nil {
			return &[]Report{}, err
		}

		reports = append(reports, report)
	}

	return &reports, nil
}

func (r *repository) GetReportById(ctx context.Context, id int) (*ReportResponse, error) {
	var report ReportResponse
	query := "SELECT * from report where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&report.Id, &report.Content, &report.Attachment, &report.Name, &report.Email, &report.Phone, &report.UserId, &report.CreatedAt)
	if err != nil {
		return &ReportResponse{}, err
	}

	report.CreatedAt, _ = util.FormatIndonesianDate(report.CreatedAt)

	return &report, nil
}

func (r *repository) GetReportByUserIdById(ctx context.Context, id int) (*ReportResponse, error) {
	var report ReportResponse
	query := "SELECT id, content, attachment FROM report where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&report.Id, &report.Content, &report.Attachment)
	if err != nil {
		return &ReportResponse{}, err
	}

	return &report, nil
}

func (r *repository) CreateReport(ctx context.Context, report *Report) error {
	fmt.Println(report.Phone)
	query := "insert into report (name, email, phone, content , attachment, user_id) values (?, ?, ?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, report.Name, report.Email, report.Phone, report.Content, report.Attachment, report.UserId)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteReport(ctx context.Context, id int) error {
	query := "delete from report where id = ?"
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

func (r *repository) GetReports(ctx context.Context) (*[]ReportResponse, error) {
	var reports []ReportResponse
	query := "select * from report"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return &[]ReportResponse{}, err
	}

	for rows.Next() {
		var report ReportResponse
		err := rows.Scan(&report.Id, &report.Content, &report.Attachment, &report.Name, &report.Email, &report.Phone, &report.UserId, &report.CreatedAt)
		if err != nil {
			return &[]ReportResponse{}, err
		}

		report.CreatedAt, _ = util.FormatIndonesianDate(report.CreatedAt)

		reports = append(reports, report)
	}

	return &reports, nil
}

func (r *repository) GetAllReportCount(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM report"
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) GetReportCountByUserId(ctx context.Context, id int) (int, error) {
	query := "SELECT COUNT(*) FROM report where user_id = ?"
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

package test

import (
	"BanjirEWS/report"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateReport(t *testing.T) {
	req := report.Report{
		Content:    "Ada banjir",
		Name:       "Holden",
		UserId:     120,
		Attachment: "test.jpg",
		Email:      "h@gmail.com",
		Phone:      "124",
	}

	err := report.Repository.CreateReport(report.NewRepository(db.GetDB()), context.Background(), &req)
	assert.Nil(t, err)
}

func TestRemoveReport(t *testing.T) {
	err := report.Repository.DeleteReport(report.NewRepository(db.GetDB()), context.Background(), 95)
	assert.Nil(t, err)
}

func TestGetReportById(t *testing.T) {
	r, err := report.Repository.GetReportById(report.NewRepository(db.GetDB()), context.Background(), 95)
	assert.Nil(t, err)
	assert.Equal(t, 95, r.Id)
}

func TestGetReports(t *testing.T) {
	rs, err := report.Repository.GetReports(report.NewRepository(db.GetDB()), context.Background())
	assert.Nil(t, err)
	assert.NotEmpty(t, rs)
}

func TestGetReportByUserIdById(t *testing.T) {
	r, err := report.Repository.GetReportByUserIdById(report.NewRepository(db.GetDB()), context.Background(), 95)
	assert.Nil(t, err)
	assert.Equal(t, 95, r.Id)
}

func TestGetAllReportCount(t *testing.T) {
	count, err := report.Repository.GetAllReportCount(report.NewRepository(db.GetDB()), context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

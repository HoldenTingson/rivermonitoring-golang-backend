package test

import (
	"BanjirEWS/history"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHistoryByRiverIdByTime(t *testing.T) {
	histories, err := history.Repository.GetHistoryByRiverIdByTime(history.NewRepository(db.GetDB()), context.Background(), "R002")
	assert.Nil(t, err)
	assert.Nil(t, *histories)
}

func TestGetHistoryByRiverId(t *testing.T) {
	histories, err := history.Repository.GetHistoryByRiverId(history.NewRepository(db.GetDB()), context.Background(), "R002")
	assert.Nil(t, err)
	assert.NotEmpty(t, histories)
}

func TestGetHistoryById(t *testing.T) {
	h, err := history.Repository.GetHistoryById(history.NewRepository(db.GetDB()), context.Background(), 135753)
	assert.Nil(t, err)
	assert.NotNil(t, h)
}

func TestGetHistoryCountByRiverId(t *testing.T) {
	count, err := history.Repository.GetHistoryCountByRiverId(history.NewRepository(db.GetDB()), context.Background(), "R001")
	assert.Nil(t, err)
	assert.Greater(t, count, 0)
}
func TestDeleteAllHistory(t *testing.T) {
	err := history.Repository.DeleteAllHistory(history.NewRepository(db.GetDB()), context.Background())
	assert.Nil(t, err)
}

func TestDeleteHistoryByRiverId(t *testing.T) {
	err := history.Repository.DeleteHistoryByRiverId(history.NewRepository(db.GetDB()), context.Background(), "R003")
	assert.Nil(t, err)
}

func TestDeleteHistoryById(t *testing.T) {
	err := history.Repository.DeleteHistoryById(history.NewRepository(db.GetDB()), context.Background(), 135755)
	assert.Nil(t, err)
}

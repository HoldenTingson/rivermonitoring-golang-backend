package history

import (
	"context"
	"time"
)

type History struct {
	Id        int     `db:"id" json:"id"`
	Height    float64 `db:"height" json:"height"`
	Status    string  `db:"status" json:"status"`
	Timestamp string  `db:"timestamp" json:"timestamp"`
	RiverId   string  `db:"river_id" json:"river_id"`
}

type HistoryResponse struct {
	Id        int     `db:"id" json:"id"`
	Height    float64 `db:"height" json:"height"`
	Status    string  `db:"status" json:"status"`
	Timestamp string  `db:"timestamp" json:"timestamp"`
	RiverId   string  `db:"river_id" json:"river_id"`
}

type Repository interface {
	GetHistoryByRiverIdByTime(ctx context.Context, id string) (*[]History, error)
	GetHistoryByRiverId(ctx context.Context, id string) (*[]History, error)
	GetHistoryById(ctx context.Context, id int) (*History, error)
	GetHistoryCountByRiverId(ctx context.Context, id string) (int, error)
	DeleteAllHistoryByTime(ctx context.Context, startTime, endTime time.Time) error
	DeleteAllHistory(ctx context.Context) error
	DeleteHistoryByRiverId(ctx context.Context, id string) error
	DeleteHistoryByRiverIdByTime(ctx context.Context, id string, startTime, endTime time.Time) error
	DeleteHistoryById(ctx context.Context, id int) error
}

type Service interface {
	ViewHistoryByRiverIdByTime(ctx context.Context, id string) (*[]HistoryResponse, error)
	ViewHistoryByRiverId(ctx context.Context, id string) (*[]HistoryResponse, error)
	ViewHistoryById(ctx context.Context, id int) (*HistoryResponse, error)
	ViewHistoryCountByRiverId(ctx context.Context, id string) (int, error)
	RemoveAllHistoryByTime(ctx context.Context, startTime, endTime time.Time) error
	RemoveAllHistory(ctx context.Context) error
	RemoveHistoryByRiverId(ctx context.Context, id string) error
	RemoveHistoryByRiverIdByTime(ctx context.Context, id string, startTime, endTime time.Time) error
	RemoveHistoryById(ctx context.Context, id int) error
}

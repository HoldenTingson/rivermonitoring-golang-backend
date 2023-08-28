package history

import (
	"context"
	"database/sql"
	"errors"
	"time"
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

func (r *repository) GetHistoryByRiverIdByTime(ctx context.Context, id string) (*[]History, error) {
	var histories []History
	query := "SELECT id, height, status, DATE_FORMAT(timestamp, '%H:%i:%s') AS timestamp FROM history WHERE river_id = ? AND timestamp >= NOW() - INTERVAL 1 MINUTE ORDER BY timestamp ASC"
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return &[]History{}, err
	}

	for rows.Next() {
		var history History
		err := rows.Scan(&history.Id, &history.Height, &history.Status, &history.Timestamp)

		if err != nil {
			return &[]History{}, err
		}

		histories = append(histories, history)
	}
	return &histories, nil
}

func (r *repository) GetHistoryByRiverId(ctx context.Context, id string) (*[]History, error) {
	var histories []History
	query := "SELECT id, height, status, timestamp, river_id FROM history where river_id = ? ORDER BY timestamp DESC"
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return &[]History{}, err
	}

	for rows.Next() {
		var history History
		err := rows.Scan(&history.Id, &history.Height, &history.Status, &history.Timestamp, &history.RiverId)

		if err != nil {
			return &[]History{}, err
		}

		histories = append(histories, history)
	}
	return &histories, nil
}

func (r *repository) GetHistoryById(ctx context.Context, id int) (*History, error) {
	var histories History
	query := "SELECT * FROM history where id = ?"

	err := r.db.QueryRowContext(ctx, query, id).Scan(&histories.Id, &histories.Height, &histories.Status, &histories.Timestamp, &histories.RiverId)
	if err != nil {
		return &History{}, err
	}
	return &histories, nil
}

func (r *repository) GetHistoryCountByRiverId(ctx context.Context, id string) (int, error) {
	query := "SELECT COUNT(*) FROM history WHERE river_id = ?"
	var count int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) DeleteAllHistoryByTime(ctx context.Context, startTime, endTime time.Time) error {
	query := "DELETE FROM history WHERE timestamp >= ? AND timestamp <= ?"
	res, err := r.db.ExecContext(ctx, query, startTime, endTime)
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

func (r *repository) DeleteAllHistory(ctx context.Context) error {
	query := "DELETE FROM history"
	res, err := r.db.ExecContext(ctx, query)
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

func (r *repository) DeleteHistoryByRiverId(ctx context.Context, id string) error {
	query := "DELETE FROM history WHERE river_id = ?"
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

func (r *repository) DeleteHistoryById(ctx context.Context, id int) error {
	query := "DELETE FROM history WHERE id = ?"
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

func (r *repository) DeleteHistoryByRiverIdByTime(ctx context.Context, id string, startTime, endTime time.Time) error {
	query := "DELETE FROM history WHERE river_id = ? AND timestamp >= ? AND timestamp <= ?"
	res, err := r.db.ExecContext(ctx, query, startTime, endTime)
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

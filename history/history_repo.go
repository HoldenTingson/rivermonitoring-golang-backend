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
	var allHistories []History
	var result []History

	var oldestTimestamp time.Time
	var rawTimestamp []uint8

	checkQuery := "SELECT MIN(timestamp) FROM history WHERE river_id = ?"
	err := r.db.QueryRowContext(ctx, checkQuery, id).Scan(&rawTimestamp)
	if err != nil && err != sql.ErrNoRows {
		return &[]History{}, err
	}

	if len(rawTimestamp) > 0 {
		parsedTime, parseErr := time.Parse("2006-01-02 15:04:05", string(rawTimestamp))
		if parseErr != nil {
			return &[]History{}, parseErr
		}
		oldestTimestamp = parsedTime
	}

	now := time.Now()
	oneMinuteAgo := now.Add(-1 * time.Minute)
	var startTime time.Time
	if oldestTimestamp.IsZero() || oldestTimestamp.Before(oneMinuteAgo) {
		startTime = oneMinuteAgo
	} else {
		startTime = oldestTimestamp
	}

	nowFormatted := now.Format("2006-01-02 15:04:05")
	startTimeFormatted := startTime.Format("2006-01-02 15:04:05")

	query := `
		SELECT height, status, timestamp
		FROM history
		WHERE river_id = ?
			AND timestamp BETWEEN ? AND ?
		ORDER BY timestamp ASC
	`
	rows, err := r.db.QueryContext(ctx, query, id, startTimeFormatted, nowFormatted)
	if err != nil {
		return &[]History{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var h History

		err := rows.Scan(&h.Height, &h.Status, &h.Timestamp)

		if err != nil {
			continue
		}

		allHistories = append(allHistories, h)
	}

	var lastIncludedTime time.Time
	for _, h := range allHistories {
		t, _ := time.Parse("2006-01-02 15:04:05", h.Timestamp)

		if lastIncludedTime.IsZero() || t.Sub(lastIncludedTime) >= 5*time.Second {
			result = append(result, h)
			lastIncludedTime = t
		}

		if len(result) >= 13 {
			break
		}
	}

	return &result, nil
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

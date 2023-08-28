package news

import (
	"BanjirEWS/util"
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
	return &repository{db: db}
}

func (r *repository) GetNews(ctx context.Context, types string) (*[]News, error) {
	var news []News
	var query string
	switch types {
	case "trending":
		query = "SELECT * FROM news where category = 'trending' ORDER BY created_at DESC limit 6"
	case "latest":
		query = "SELECT * FROM news ORDER BY id DESC limit 5"
	case "other":
		query = "SELECT * FROM news where category = 'other' ORDER BY created_at DESC  limit 5"
	case "main":
		query = "SELECT * FROM news where category = 'main' ORDER BY created_at DESC limit 1"
	default:
		query = "SELECT * FROM news"
	}
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return &[]News{}, err
	}

	for rows.Next() {
		var new News
		err := rows.Scan(&new.Id, &new.Title, &new.Content, &new.Description, &new.Image, &new.Category, &new.CreatedAt)

		if err != nil {
			return &[]News{}, err
		}

		new.CreatedAt, _ = util.FormatIndonesianDate(new.CreatedAt)

		news = append(news, new)
	}
	return &news, nil

}

func (r *repository) GetNewsByID(ctx context.Context, id int) (*News, error) {
	var news News
	query := "select id, title, content, description, image, category, created_at FROM news where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&news.Id, &news.Title, &news.Content, &news.Description, &news.Image, &news.Category, &news.CreatedAt)
	if err != nil {
		return &News{}, err
	}

	news.CreatedAt, _ = util.FormatIndonesianDate(news.CreatedAt)

	return &news, nil
}

func (r *repository) CreateNews(ctx context.Context, news *News) error {
	query := "insert into news (title, content, description, image, category) values(?,?,?,?,?)"
	_, err := r.db.ExecContext(ctx, query, news.Title, news.Content, news.Description, news.Image, news.Category)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateNews(ctx context.Context, news *News, id int) error {
	query := "update news set id = ?, title = ?, content = ?, description = ?, image = ?, category = ? where id = ?"
	res, err := r.db.ExecContext(ctx, query, news.Id, news.Title, news.Content, news.Description, news.Image, news.Category, id)
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

func (r *repository) DeleteNews(ctx context.Context, id int) error {
	query := "delete from news where id = ?"
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

func (r *repository) GetAllNewsCount(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM news"
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

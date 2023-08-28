package gallery

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

func (r *repository) GetGallery(ctx context.Context) (*[]Gallery, error) {
	var galleries []Gallery
	query := "select * FROM gallery"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return &[]Gallery{}, err
	}

	for rows.Next() {
		var g Gallery
		if err := rows.Scan(&g.Id, &g.Title, &g.Image, &g.Date); err != nil {
			return nil, err
		}

		g.Date, _ = util.FormatIndonesianDate(g.Date)
		galleries = append(galleries, g)
	}

	return &galleries, nil
}

func (r *repository) GetGalleryById(ctx context.Context, id int) (*Gallery, error) {
	var gallery Gallery
	query := "select * FROM gallery where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&gallery.Id, &gallery.Title, &gallery.Image, &gallery.Date)
	if err != nil {
		return &Gallery{}, err
	}
	gallery.Date, _ = util.FormatIndonesianDate(gallery.Date)

	return &gallery, nil
}

func (r *repository) GetGalleryByIdAdmin(ctx context.Context, id int) (*Gallery, error) {
	var gallery Gallery
	query := "select * FROM gallery where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&gallery.Id, &gallery.Title, &gallery.Image, &gallery.Date)
	if err != nil {
		return &Gallery{}, err
	}

	return &gallery, nil
}

func (r *repository) CreateGallery(ctx context.Context, gallery *Gallery) error {
	query := "insert into gallery (title, image,  date) values(?,?,?)"
	_, err := r.db.ExecContext(ctx, query, gallery.Title, gallery.Image, gallery.Date)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateGallery(ctx context.Context, gallery *Gallery, id int) error {
	query := "update gallery set id = ?, title = ?, image = ?,  date = ? where id = ?"
	res, err := r.db.ExecContext(ctx, query, gallery.Id, gallery.Title, gallery.Image, gallery.Date, id)
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

func (r *repository) DeleteGallery(ctx context.Context, id int) error {
	query := "delete from gallery where id = ?"
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

func (r *repository) GetAllGalleryCount(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM gallery"
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

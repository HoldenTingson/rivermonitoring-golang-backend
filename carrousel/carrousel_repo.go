package carrousel

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

func (r *repository) GetCarrousel(ctx context.Context) (*[]Carrousel, error) {
	var carrousel []Carrousel
	query := "select * from carrousel"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return &[]Carrousel{}, err
	}

	for rows.Next() {
		var new Carrousel
		err := rows.Scan(&new.Id, &new.Title, &new.Desc, &new.Image, &new.Date)

		if err != nil {
			return &[]Carrousel{}, err
		}

		// Parse the date string into a time.Time value
		new.Date, _ = util.FormatIndonesianDate(new.Date)

		carrousel = append(carrousel, new)
	}
	return &carrousel, nil

}

func (r *repository) GetCarrouselByID(ctx context.Context, id int) (*Carrousel, error) {
	var carrousel Carrousel
	query := "select id, title, description, image, date FROM carrousel where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&carrousel.Id, &carrousel.Title, &carrousel.Desc, &carrousel.Image, &carrousel.Date)
	carrousel.Date, _ = util.FormatIndonesianDate(carrousel.Date)
	if err != nil {
		return &Carrousel{}, err
	}
	return &carrousel, nil
}

func (r *repository) GetCarrouselByIDAdmin(ctx context.Context, id int) (*Carrousel, error) {
	var carrousel Carrousel
	query := "select id, title, description, image, date FROM carrousel where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&carrousel.Id, &carrousel.Title, &carrousel.Desc, &carrousel.Image, &carrousel.Date)
	if err != nil {
		return &Carrousel{}, err
	}
	return &carrousel, nil
}

func (r *repository) CreateCarrousel(ctx context.Context, carrousel *Carrousel) error {
	query := "insert into carrousel (title, description, image, date) values(?,?,?,?)"
	_, err := r.db.ExecContext(ctx, query, carrousel.Title, carrousel.Desc, carrousel.Image, carrousel.Date)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateCarrousel(ctx context.Context, carrousel *Carrousel, id int) error {
	query := "update carrousel set title = ?, description = ?, image = ?, date = ? where id = ?"
	res, err := r.db.ExecContext(ctx, query, carrousel.Title, carrousel.Desc, carrousel.Image, carrousel.Date, id)
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

func (r *repository) DeleteCarrousel(ctx context.Context, id int) error {
	query := "delete from carrousel where id = ?"
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

func (r *repository) GetAllCarrouselCount(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM carrousel"
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

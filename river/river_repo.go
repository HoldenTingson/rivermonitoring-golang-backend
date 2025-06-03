package river

import (
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

func (r *repository) UpdateRiver(ctx context.Context, update *UpdateRiver) error {

	query := "UPDATE river SET height = ?, status = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, update.Height, update.Status, update.Id)
	if err != nil {
		return err
	}

	insertQuery := "INSERT INTO history (height, status, river_id) VALUES (?, ?, ?)"
	_, err = r.db.ExecContext(ctx, insertQuery, update.Height, update.Status, update.Id)
	if err != nil {
		return err
	}

	fmt.Println("River updated:", update)
	return nil
}

func (r *repository) GetRiver(ctx context.Context) (*[]River, error) {
	var rivers []River
	query := "select * from river"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return &[]River{}, err
	}

	for rows.Next() {
		var river River
		err := rows.Scan(&river.Id, &river.Latitude, &river.Longitude, &river.Location, &river.Height, &river.Status)

		if err != nil {
			return &[]River{}, err
		}

		rivers = append(rivers, river)
	}
	return &rivers, nil
}

func (r *repository) GetRiverById(ctx context.Context, id string) (*River, error) {
	var river River
	query := "select * FROM river where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&river.Id, &river.Latitude, &river.Longitude, &river.Location, &river.Height, &river.Status)
	if err != nil {
		return &River{}, err
	}

	return &river, nil
}

func (r *repository) GetRiverByStatus(ctx context.Context, status string) (*[]River, error) {
	var rivers []River
	query := "select * from river where status = ?"
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return &[]River{}, err
	}

	for rows.Next() {
		var river River
		err := rows.Scan(&river.Id, &river.Latitude, &river.Longitude, &river.Location, &river.Height, &river.Status)

		if err != nil {
			return &[]River{}, err
		}

		rivers = append(rivers, river)
	}
	return &rivers, nil
}

func (r *repository) GetRiverId(ctx context.Context) ([]string, error) {
	query := "select id from river"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var riverIDs []string
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		riverIDs = append(riverIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return riverIDs, nil
}

func (r *repository) FilterRiver(ctx context.Context, sortBy string) (*[]River, error) {

	switch sortBy {
	case "locationAsc":
		sortBy = "location ASC"
	case "locationDesc":
		sortBy = "location DESC"
	case "heightAsc":
		sortBy = "height ASC"
	case "heightDesc":
		sortBy = "height DESC"
	case "latitudeAsc":
		sortBy = "latitude ASC"
	case "latitudeDesc":
		sortBy = "latitude DESC"
	case "longitudeAsc":
		sortBy = "longitude ASC"
	case "longitudeDesc":
		sortBy = "longitude DESC"
	case "statusAsc":
		sortBy = "status ASC"
	case "statusDesc":
		sortBy = "status DESC"
	}

	query := fmt.Sprintf("SELECT id, latitude, longitude, location, height, status FROM river ORDER BY %s", sortBy)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rivers []River
	for rows.Next() {
		var river River
		err := rows.Scan(&river.Id, &river.Latitude, &river.Longitude, &river.Location, &river.Height, &river.Status)
		if err != nil {
			return nil, err
		}
		rivers = append(rivers, river)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &rivers, nil
}

func (r *repository) FindRiver(ctx context.Context, searchTerm string) (*River, error) {
	var river River
	query := "SELECT * FROM river WHERE location LIKE ? limit 1"
	searchTerm = "%" + searchTerm + "%"
	err := r.db.QueryRowContext(ctx, query, searchTerm).Scan(&river.Id, &river.Latitude, &river.Longitude, &river.Location, &river.Height, &river.Status)
	if err != nil {
		return &River{}, err
	}

	return &river, nil
}

func (r *repository) GetAllRiverCount(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM river"
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) CreateRiver(ctx context.Context, river *River) error {
	query := "insert into river (id, latitude, longitude, location) values(?,?,?,?)"
	_, err := r.db.ExecContext(ctx, query, river.Id, river.Latitude, river.Longitude, river.Location)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteRiver(ctx context.Context, id string) error {
	query := "delete from river where id = ?"
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

func (r *repository) UpdateRiverDetail(ctx context.Context, river *River, id string) error {
	query := "update river set id = ?, latitude = ?, longitude = ?, location = ?, height = ?, status = ? where id = ?"
	res, err := r.db.ExecContext(ctx, query, river.Id, river.Latitude, river.Longitude, river.Location, river.Height, river.Status, id)
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

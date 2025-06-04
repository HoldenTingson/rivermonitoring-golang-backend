package admin

import (
	"BanjirEWS/util"
	"context"
	"database/sql"
	"errors"
	"log"
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

func (r *repository) CreateAdmin(ctx context.Context, admin *Admin) (*Admin, error) {
	query := "insert into admin (username , password) values (?, ?)"
	res, err := r.db.ExecContext(ctx, query, admin.Username, admin.Password)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	admin.Id = int(id)

	return admin, nil
}

func (r *repository) GetAdminByUsername(ctx context.Context, username string) (*Admin, error) {
	var admin Admin
	query := "SELECT * FROM admin WHERE username = ?"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&admin.Id, &admin.Username, &admin.Password, &admin.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *repository) GetAdmin(ctx context.Context, token string) (*Admin, error) {
	var admin Admin
	query := "select * from admin where id = ?"
	err := r.db.QueryRowContext(ctx, query, token).Scan(&admin.Id, &admin.Username, &admin.Password, &admin.CreatedAt)
	if err != nil {
		return &Admin{}, nil
	}

	return &admin, nil
}

func (r *repository) GetAllUserCount(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM user"
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) GetAdmins(ctx context.Context) (*[]Admin, error) {
	var admins []Admin
	query := "select id, username, created_at from admin"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var admin Admin
		if err := rows.Scan(&admin.Id, &admin.Username, &admin.CreatedAt); err != nil {
			return nil, err
		}
		admin.CreatedAt, _ = util.FormatIndonesianTimezone(admin.CreatedAt)
		admins = append(admins, admin)
	}

	return &admins, nil

}

func (r *repository) GetAdminById(ctx context.Context, id int) (*Admin, error) {
	var admin Admin
	query := "select id, username, created_at FROM admin where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&admin.Id, &admin.Username, &admin.CreatedAt)
	if err != nil {
		return &Admin{}, err
	}

	admin.CreatedAt, _ = util.FormatIndonesianTimezone(admin.CreatedAt)
	log.Println(admin.CreatedAt)

	return &admin, nil
}

func (r *repository) DeleteAdmin(ctx context.Context, id int) error {
	query := "delete from admin where id = ?"
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

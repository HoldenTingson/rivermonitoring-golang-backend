package admin

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

func (r *repository) GetAdminByUsername(ctx context.Context, username string, password string) (*Admin, error) {
	var admin Admin
	query := "SELECT * FROM admin WHERE username = ?"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&admin.Id, &admin.Username, &admin.Password)
	if err != nil {
		return nil, err
	}

	// Check if the provided password matches the stored password
	if admin.Password != password {
		return nil, errors.New("incorrect password")
	}

	return &admin, nil
}

func (r *repository) GetAdmin(ctx context.Context, token string) (*Admin, error) {
	var admin Admin
	query := "select * from admin where id = ?"
	err := r.db.QueryRowContext(ctx, query, token).Scan(&admin.Id, &admin.Username, &admin.Password)
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

func (r *repository) GetUser(ctx context.Context) (*[]User, error) {
	var users []User
	query := "select * FROM user"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return &[]User{}, err
	}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Language, &user.Profile, &user.CreatedAt, &user.ChangedAt); err != nil {
			return nil, err
		}

		user.CreatedAt, _ = util.FormatIndonesianDate(user.CreatedAt)
		user.ChangedAt, _ = util.FormatIndonesianDate(user.ChangedAt)

		users = append(users, user)
	}

	return &users, nil
}

func (r *repository) GetUserById(ctx context.Context, id int) (*User, error) {
	var user User
	query := "select * FROM user where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Language, &user.Profile, &user.CreatedAt, &user.ChangedAt)
	if err != nil {
		return &User{}, err
	}

	user.CreatedAt, _ = util.FormatIndonesianDate(user.CreatedAt)
	user.ChangedAt, _ = util.FormatIndonesianDate(user.ChangedAt)

	return &user, nil
}

func (r *repository) GetUserByIdAdmin(ctx context.Context, id int) (*User, error) {
	var user User
	query := "select * FROM user where id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Language, &user.Profile, &user.CreatedAt, &user.ChangedAt)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	query := "insert into user (username, password, email, phone, language, profile) values(?,?,?,?,?,?)"
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Password, user.Email, user.Phone, user.Language, user.Profile)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User, id int) error {
	query := "update user set username = ?, password = ?, email = ?, phone = ?, language = ?, profile = ? where id = ? "
	res, err := r.db.ExecContext(ctx, query, user.Username, user.Password, user.Email, user.Phone, user.Language, user.Profile, id)
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

func (r *repository) DeleteUser(ctx context.Context, id int) error {
	query := "delete from user where id = ?"
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

package user

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

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	query := "insert into user (username , password, email) values (?, ?, ?)"
	res, err := r.db.ExecContext(ctx, query, user.Username, user.Password, user.Email)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.Id = id

	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}
	query := "select id, username, password, email from user where email = ?"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.Id, &u.Username, &u.Password, &u.Email)
	if err != nil {
		return &User{}, nil
	}

	return &u, nil
}

func (r *repository) GetUser(ctx context.Context, token string) (*User, error) {
	u := User{}
	query := "select id, username, password, email, phone, language, profile from user where id = ?"
	err := r.db.QueryRowContext(ctx, query, token).Scan(&u.Id, &u.Username, &u.Password, &u.Email, &u.Phone, &u.Language, &u.Profile)
	if err != nil {

		return &User{}, nil
	}
	return &u, nil
}

func (r *repository) UpdateProfile(ctx context.Context, user *User, id int, update string) error {
	var query string
	var args []interface{}
	switch update {
	case "account":
		query = "update user set username = ?, email = ?, phone = ?, language = ?, changed_at = ?  where id = ?"
		args = []interface{}{user.Username, user.Email, user.Phone, user.Language, time.Now(), id}
	case "password":
		query = "update user set password = ?, changed_at = ? where id = ?"
		args = []interface{}{user.Password, time.Now(), id}
	case "profile":
		query = "update user set profile = ?, changed_at = ? where id = ?"
		args = []interface{}{user.Profile, time.Now(), id}
	}
	res, err := r.db.ExecContext(ctx, query, args...)
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

func (r *repository) UpdatePassword(ctx context.Context, email string, newPassword string) error {
	query := "UPDATE user SET password = ? WHERE email = ?"
	_, err := r.db.ExecContext(ctx, query, newPassword, email)
	if err != nil {
		return err
	}
	return nil
}

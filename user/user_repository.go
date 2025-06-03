package user

import (
	"BanjirEWS/util"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
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

func (r *repository) RegisterUser(ctx context.Context, user *User) (*User, error) {
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

	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	emailVerPassword := string(emailVerRandRune)

	var emailVerPasswordHash []byte
	emailVerPasswordHash, err = bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	emailVerPasswordHashStr := string(emailVerPasswordHash)

	timeout := time.Now().Local().AddDate(0, 0, 1)

	query = "insert into email_ver_hash (username , email, hash, timeout) values (?, ?, ?, ?)"
	res, err = r.db.ExecContext(ctx, query, user.Username, user.Email, emailVerPasswordHashStr, timeout)
	if err != nil {
		return nil, err
	}

	if rows, err := res.RowsAffected(); rows != 1 {
		return nil, err
	}

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
	query := "select id, username, password, email, phone, profile from user where id = ?"
	err := r.db.QueryRowContext(ctx, query, token).Scan(&u.Id, &u.Username, &u.Password, &u.Email, &u.Phone, &u.Profile)
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
		query = "update user set username = ?, email = ?, phone = ?, changed_at = ?  where id = ?"
		args = []interface{}{user.Username, user.Email, user.Phone, time.Now(), id}
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

func (r *repository) CheckPassword(ctx context.Context, email string, currentPassword string) (bool, error) {
	query := "select password from user where email = ?"
	var hashedPassword string
	err := r.db.QueryRowContext(ctx, query, email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("user not found")
		}
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword))
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (r *repository) SendEmail(ctx context.Context, email string, hash string) (string, error) {
	var username string
	query := "select email, username from user where email = ?"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&email, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return username, errors.New("user not found")
		}
		return username, err
	}

	now := time.Now()
	timeout := now.Add(time.Minute * 45)

	query = "update email_ver_hash set hash = ?, timeout = ? where email = ?"
	_, err = r.db.ExecContext(ctx, query, hash, timeout, email)
	if err != nil {
		return username, err
	}

	return username, nil
}

func (r *repository) CreatePassword(ctx context.Context, username string, hash string, password string) error {
	query := "select hash, timeout from email_ver_hash where username = ?"
	var dbHash string
	var timeout string

	err := r.db.QueryRowContext(ctx, query, username).Scan(&dbHash, &timeout)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}

		return err
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeout)
	if err != nil {
		return err
	}
	currentTime := time.Now()

	if currentTime.After(parsedTime) {
		return errors.New("forgot password link already expired")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(hash))
	if err != nil {
		return err
	}

	query = "update user set password = ? where username = ?"
	res, err := r.db.ExecContext(ctx, query, password, username)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	if rows != 1 {
		return errors.New("no rows affected")
	}

	return nil
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	query := "insert into user (username, password, email, phone, profile) values(?,?,?,?,?)"
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Password, user.Email, user.Phone, user.Profile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	emailVerPassword := string(emailVerRandRune)

	var emailVerPasswordHash []byte
	emailVerPasswordHash, err = bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return err
	}
	emailVerPasswordHashStr := string(emailVerPasswordHash)

	timeout := time.Now().Local().AddDate(0, 0, 1)

	query = "insert into email_ver_hash (username , email, hash, timeout) values (?, ?, ?, ?)"
	res, err := r.db.ExecContext(ctx, query, user.Username, user.Email, emailVerPasswordHashStr, timeout)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if rows, err := res.RowsAffected(); rows != 1 {
		fmt.Println(err)
		return err
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

func (r *repository) GetUsers(ctx context.Context) (*[]User, error) {
	var users []User
	query := "select * FROM user"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return &[]User{}, err
	}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Profile, &user.CreatedAt, &user.ChangedAt); err != nil {
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
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Profile, &user.CreatedAt, &user.ChangedAt)
	if err != nil {
		return &User{}, err
	}

	user.CreatedAt, _ = util.FormatIndonesianDate(user.CreatedAt)
	user.ChangedAt, _ = util.FormatIndonesianDate(user.ChangedAt)

	return &user, nil
}

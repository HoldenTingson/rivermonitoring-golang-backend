package faq

import (
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

func (r *repository) GetFaq(ctx context.Context) (*[]Faq, error) {
	var faq []Faq
	query := "select * FROM faq"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return &[]Faq{}, err
	}

	for rows.Next() {
		var f Faq
		if err := rows.Scan(&f.Id, &f.Category, &f.Question, &f.Answer); err != nil {
			return nil, err
		}
		faq = append(faq, f)
	}

	return &faq, nil
}

func (r *repository) CreateFaq(ctx context.Context, faq *Faq) error {
	query := "insert into faq (category, question, answer) values(?,?,?)"
	_, err := r.db.ExecContext(ctx, query, faq.Category, faq.Question, faq.Answer)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateFaq(ctx context.Context, faq *Faq, id int) error {
	query := "update faq set category = ?, question = ?, answer = ? where id = ?"
	res, err := r.db.ExecContext(ctx, query, faq.Category, faq.Question, faq.Answer, id)
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

func (r *repository) DeleteFaq(ctx context.Context, id int) error {
	query := "delete from faq where id = ?"
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

func (r *repository) FindFaq(ctx context.Context, question string) (*[]Faq, error) {
	query := "SELECT id, question, answer FROM faq WHERE question LIKE ?"
	question = "%" + question + "%"
	rows, err := r.db.QueryContext(ctx, query, question)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faq []Faq
	for rows.Next() {
		var f Faq
		err := rows.Scan(&f.Id, &f.Question, &f.Answer)
		if err != nil {
			return nil, err
		}
		faq = append(faq, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &faq, nil
}

func (r *repository) GetCategory(ctx context.Context) (*[]Faq, error) {
	query := "SELECT DISTINCT category FROM faq"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faq []Faq
	for rows.Next() {
		var f Faq
		err := rows.Scan(&f.Category)
		if err != nil {
			return nil, err
		}
		faq = append(faq, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &faq, nil
}

func (r *repository) GetQa(ctx context.Context, category string) (*[]Faq, error) {
	query := "SELECT id, question, answer FROM faq where category = ?"
	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faq []Faq
	for rows.Next() {
		var f Faq
		err := rows.Scan(&f.Id, &f.Question, &f.Answer)
		if err != nil {
			return nil, err
		}
		faq = append(faq, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &faq, nil
}

func (r *repository) GetAllFaqCount(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM faq"
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

package sqlexpr

import (
	"context"
	"database/sql"
)

// Executor is compatible with *sql.DB, *sql.Tx (or an arbitrary middleman)
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func Exec(ctx context.Context, ex Executor, expr Expr) (sql.Result, error) {
	query, args := Build(expr)
	return ex.ExecContext(ctx, query, args...)
}

func Query(ctx context.Context, ex Executor, expr Expr) (*sql.Rows, error) {
	query, args := Build(expr)
	return ex.QueryContext(ctx, query, args...)
}

func QueryRow(ctx context.Context, ex Executor, expr Expr) *sql.Row {
	query, args := Build(expr)
	return ex.QueryRowContext(ctx, query, args...)
}

func (s *Select) Query(ctx context.Context, ex Executor) (*sql.Rows, error) {
	return Query(ctx, ex, s)
}

func (s *Select) QueryRow(ctx context.Context, ex Executor) *sql.Row {
	return QueryRow(ctx, ex, s)
}

func (s *Insert) Exec(ctx context.Context, ex Executor) (sql.Result, error) {
	return Exec(ctx, ex, s)
}

func (s *Insert) Query(ctx context.Context, ex Executor) (*sql.Rows, error) {
	return Query(ctx, ex, s)
}

func (s *Insert) QueryRow(ctx context.Context, ex Executor) *sql.Row {
	return QueryRow(ctx, ex, s)
}

func (s *Update) Exec(ctx context.Context, ex Executor) (sql.Result, error) {
	return Exec(ctx, ex, s)
}

func (s *Update) Query(ctx context.Context, ex Executor) (*sql.Rows, error) {
	return Query(ctx, ex, s)
}

func (s *Update) QueryRow(ctx context.Context, ex Executor) *sql.Row {
	return QueryRow(ctx, ex, s)
}

func (s *Delete) Exec(ctx context.Context, ex Executor) (sql.Result, error) {
	return Exec(ctx, ex, s)
}

func (s *Delete) Query(ctx context.Context, ex Executor) (*sql.Rows, error) {
	return Query(ctx, ex, s)
}

func (s *Delete) QueryRow(ctx context.Context, ex Executor) *sql.Row {
	return QueryRow(ctx, ex, s)
}

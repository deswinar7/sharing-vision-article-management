package article

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("article not found")

type Repository interface {
	Create(ctx context.Context, input CreateInput) (Article, error)
	List(ctx context.Context, filter ListFilter) ([]Article, error)
	GetByID(ctx context.Context, id int64) (Article, error)
	Update(ctx context.Context, article Article) (Article, error)
	Delete(ctx context.Context, id int64) error
}

type SQLRepository struct{ db *sqlx.DB }

func NewSQLRepository(db *sqlx.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) Create(ctx context.Context, input CreateInput) (Article, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO posts (title, content, category, status) VALUES (?, ?, ?, ?)`,
		input.Title, input.Content, input.Category, input.Status,
	)
	if err != nil {
		return Article{}, fmt.Errorf("insert article: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Article{}, fmt.Errorf("read inserted article id: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]Article, error) {
	query := `SELECT id, title, content, category, created_date, updated_date, status FROM posts`
	args := make([]any, 0, 3)
	if filter.Status != nil {
		query += ` WHERE status = ?`
		args = append(args, *filter.Status)
	}
	query += ` ORDER BY created_date DESC, id DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Limit, filter.Offset)

	articles := make([]Article, 0)
	if err := r.db.SelectContext(ctx, &articles, query, args...); err != nil {
		return nil, fmt.Errorf("list articles: %w", err)
	}
	return articles, nil
}

func (r *SQLRepository) GetByID(ctx context.Context, id int64) (Article, error) {
	var result Article
	err := r.db.GetContext(ctx, &result,
		`SELECT id, title, content, category, created_date, updated_date, status FROM posts WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Article{}, ErrNotFound
	}
	if err != nil {
		return Article{}, fmt.Errorf("get article: %w", err)
	}
	return result, nil
}

func (r *SQLRepository) Update(ctx context.Context, article Article) (Article, error) {
	result, err := r.db.ExecContext(ctx,
		`UPDATE posts SET title = ?, content = ?, category = ?, status = ? WHERE id = ?`,
		article.Title, article.Content, article.Category, article.Status, article.ID,
	)
	if err != nil {
		return Article{}, fmt.Errorf("update article: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return Article{}, fmt.Errorf("read updated row count: %w", err)
	}
	if affected == 0 {
		return Article{}, ErrNotFound
	}
	return r.GetByID(ctx, article.ID)
}

func (r *SQLRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete article: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read deleted row count: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

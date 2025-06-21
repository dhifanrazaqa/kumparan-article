package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrArticleNotFound = errors.New("article not found")

type ArticleRepository interface {
	Create(ctx context.Context, article *models.Article) error
	FindByID(ctx context.Context, id string) (*models.Article, error)
	FindAll(ctx context.Context, params models.ListArticlesParams) ([]models.Article, error)
	Update(ctx context.Context, article *models.Article) error
	Delete(ctx context.Context, id string) error
	CountAll(ctx context.Context, params models.ListArticlesParams) (int64, error)
}

type pgxArticleRepo struct {
	pool *pgxpool.Pool
}

func NewPgxArticleRepo(pool *pgxpool.Pool) ArticleRepository {
	return &pgxArticleRepo{pool: pool}
}

func (r *pgxArticleRepo) Create(ctx context.Context, article *models.Article) error {
	query := `INSERT INTO articles (title, body, author_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	row := r.pool.QueryRow(ctx, query, article.Title, article.Body, article.AuthorID)
	err := row.Scan(&article.ID, &article.CreatedAt, &article.UpdatedAt)
	return err
}

func (r *pgxArticleRepo) FindByID(ctx context.Context, id string) (*models.Article, error) {
	query := `
		SELECT
			a.id, a.title, a.body, a.author_id, a.created_at, a.updated_at,
			u.username as author_username, u.name as author_name, u.created_at as author_created_at, u.updated_at as author_updated_at
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE a.id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var article models.Article
	var author models.UserResponse
	err := row.Scan(
		&article.ID, &article.Title, &article.Body, &article.AuthorID, &article.CreatedAt, &article.UpdatedAt,
		&author.Username, &author.Name, &author.CreatedAt, &author.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}
	author.ID = article.AuthorID
	article.Author = &author
	return &article, nil
}

func (r *pgxArticleRepo) FindAll(ctx context.Context, params models.ListArticlesParams) ([]models.Article, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT
			a.id, a.title, a.body, a.author_id, a.created_at, a.updated_at,
			u.username, u.name, u.created_at, u.updated_at
		FROM articles a
		JOIN users u ON a.author_id = u.id
	`)

	var args []interface{}
	var conditions []string

	if params.Author != "" {
		args = append(args, params.Author)
		conditions = append(conditions, fmt.Sprintf("LOWER(u.name) = LOWER($%d)", len(args)))
	}
	if params.Query != "" {
		searchQuery := strings.Join(strings.Fields(params.Query), " & ")
		args = append(args, searchQuery)
		conditions = append(conditions, fmt.Sprintf("a.search_vector @@ to_tsquery('english', $%d)", len(args)))
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY a.created_at DESC")
	args = append(args, params.Limit)
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", len(args)))
	args = append(args, params.Offset)
	queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%d", len(args)))

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan query artikel: %w", err)
	}
	defer rows.Close()

	articles := make([]models.Article, 0)
	for rows.Next() {
		var article models.Article
		var author models.UserResponse
		err := rows.Scan(
			&article.ID, &article.Title, &article.Body, &article.AuthorID, &article.CreatedAt, &article.UpdatedAt,
			&author.Username, &author.Name, &author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article row: %w", err)
		}
		author.ID = article.AuthorID
		article.Author = &author
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return articles, nil
}

func (r *pgxArticleRepo) CountAll(ctx context.Context, params models.ListArticlesParams) (int64, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT COUNT(*) FROM articles a")

	var args []interface{}
	var conditions []string

	if params.Author != "" || params.Query != "" {
		queryBuilder.WriteString(" JOIN users u ON a.author_id = u.id")
	}

	if params.Author != "" {
		args = append(args, params.Author)
		conditions = append(conditions, fmt.Sprintf("LOWER(u.name) = LOWER($%d)", len(args)))
	}
	if params.Query != "" {
		searchQuery := strings.Join(strings.Fields(params.Query), " & ")
		args = append(args, searchQuery)
		conditions = append(conditions, fmt.Sprintf("a.search_vector @@ to_tsquery('english', $%d)", len(args)))
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	var count int64
	err := r.pool.QueryRow(ctx, queryBuilder.String(), args...).Scan(&count)
	return count, err
}

func (r *pgxArticleRepo) Update(ctx context.Context, article *models.Article) error {
	query := `UPDATE articles SET title = $1, body = $2 WHERE id = $3 RETURNING updated_at`
	row := r.pool.QueryRow(ctx, query, article.Title, article.Body, article.ID)
	err := row.Scan(&article.UpdatedAt)
	return err
}

func (r *pgxArticleRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM articles WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrArticleNotFound
	}
	return nil
}

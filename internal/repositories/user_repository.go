package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

type pgxUserRepo struct {
	pool *pgxpool.Pool
}

func NewPgxUserRepo(pool *pgxpool.Pool) UserRepository {
	return &pgxUserRepo{pool: pool}
}

func (r *pgxUserRepo) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, hashed_password) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	row := r.pool.QueryRow(ctx, query, user.Username, user.HashedPassword)
	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *pgxUserRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, username, hashed_password, created_at, updated_at FROM users WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *pgxUserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, hashed_password, created_at, updated_at FROM users WHERE username = $1`
	row := r.pool.QueryRow(ctx, query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *pgxUserRepo) FindAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, username, hashed_password, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *pgxUserRepo) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET username = $1, hashed_password = $2 WHERE id = $3 RETURNING updated_at`
	row := r.pool.QueryRow(ctx, query, user.Username, user.HashedPassword, user.ID)
	err := row.Scan(&user.UpdatedAt)
	return err
}

func (r *pgxUserRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/online-shop/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password, first_name, name, last_name, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		user.Email, user.Password, user.FirstName, user.FirstName, user.LastName, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	query := `SELECT id, email, password, first_name, last_name, role, created_at, updated_at
	           FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName,
		&u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	query := `SELECT id, email, password, first_name, last_name, role, created_at, updated_at
	           FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName,
		&u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

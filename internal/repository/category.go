package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/online-shop/internal/models"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, c *models.Category) error {
	query := `INSERT INTO categories (name, slug, parent_id)
	          VALUES ($1, $2, $3)
	          RETURNING id, created_at`

	return r.db.QueryRow(ctx, query, c.Name, c.Slug, c.ParentID).Scan(&c.ID, &c.CreatedAt)
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (*models.Category, error) {
	var c models.Category
	query := `SELECT id, name, slug, parent_id, created_at FROM categories WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.Name, &c.Slug, &c.ParentID, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) List(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, slug, parent_id, created_at FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.ParentID, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *models.Category) error {
	query := `UPDATE categories SET name=$1, slug=$2, parent_id=$3 WHERE id=$4`
	_, err := r.db.Exec(ctx, query, c.Name, c.Slug, c.ParentID, c.ID)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM categories WHERE id = $1", id)
	return err
}

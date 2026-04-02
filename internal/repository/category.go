package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/online-shop/internal/models"
)

type CategoryRepository struct {
	db       *pgxpool.Pool
	adminURL string
}

func NewCategoryRepository(db *pgxpool.Pool, adminURL string) *CategoryRepository {
	return &CategoryRepository{db: db, adminURL: adminURL}
}

func (r *CategoryRepository) Create(ctx context.Context, c *models.Category) error {
	query := `INSERT INTO categories (name, slug, parent_id)
	          VALUES ($1, $2, $3)
	          RETURNING id, created_at`

	return r.db.QueryRow(ctx, query, c.Name, c.Slug, c.ParentID).Scan(&c.ID, &c.CreatedAt)
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (*models.Category, error) {
	var c models.Category
	var mediaID *int64
	var mediaFileName *string

	query := `SELECT c.id, c.name, c.slug, c.parent_id, c.created_at,
	                 m.id, m.file_name
	          FROM categories c
	          LEFT JOIN media m ON m.model_id = c.id
	              AND m.model_type = 'App\Models\Category'
	              AND m.collection_name = 'preview'
	          WHERE c.id = $1
	          LIMIT 1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Slug, &c.ParentID, &c.CreatedAt,
		&mediaID, &mediaFileName,
	)
	if err != nil {
		return nil, err
	}

	if mediaID != nil && mediaFileName != nil {
		url := fmt.Sprintf("%s/storage/%d/%s", r.adminURL, *mediaID, *mediaFileName)
		c.Preview = &url
	}

	return &c, nil
}

func (r *CategoryRepository) List(ctx context.Context) ([]models.Category, error) {
	query := `SELECT c.id, c.name, c.slug, c.parent_id, c.created_at,
	                 m.id, m.file_name
	          FROM categories c
	          LEFT JOIN media m ON m.model_id = c.id
	              AND m.model_type = 'App\Models\Category'
	              AND m.collection_name = 'preview'
	          ORDER BY c.name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		var mediaID *int64
		var mediaFileName *string

		if err := rows.Scan(
			&c.ID, &c.Name, &c.Slug, &c.ParentID, &c.CreatedAt,
			&mediaID, &mediaFileName,
		); err != nil {
			return nil, err
		}

		if mediaID != nil && mediaFileName != nil {
			url := fmt.Sprintf("%s/storage/%d/%s", r.adminURL, *mediaID, *mediaFileName)
			c.Preview = &url
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

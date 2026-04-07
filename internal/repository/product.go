package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/online-shop/internal/models"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p *models.Product) error {
	query := `
		INSERT INTO products (name, slug, description, price, stock, category_id, image_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		p.Name, p.Slug, p.Description, p.Price, p.Stock,
		p.CategoryID, p.ImageURL, p.IsActive,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	var p models.Product
	query := `SELECT p.id, p.name, p.sku, p.slug, p.price, p.brand_id,b.name as brand_name, p.category_id,
	           p.content, p.created_at, p.description, p.old_price, p.is_soon, p.is_featured
	    FROM products p
		left join brands b on b.id=p.brand_id
	    WHERE p.id = $1 AND p.is_active = true AND p.deleted_at IS NULL`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Sku, &p.Slug, &p.Price, &p.BrandId, &p.BrandName, &p.CategoryID,
		&p.Content, &p.CreatedAt, &p.Description, &p.OldPrice, &p.IsSoon, &p.IsFeatured,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProductRepository) List(ctx context.Context, f models.ProductFilter) ([]models.Product, int, error) {
	var (
		conditions []string
		args       []interface{}
		argIdx     = 1
	)

	if f.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("p.category_id = $%d", argIdx))
		args = append(args, *f.CategoryID)
		argIdx++
	}
	if f.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("p.price >= $%d", argIdx))
		args = append(args, *f.MinPrice)
		argIdx++
	}
	if f.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("p.price <= $%d", argIdx))
		args = append(args, *f.MaxPrice)
		argIdx++
	}
	if f.Search != nil {
		conditions = append(conditions, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+*f.Search+"%")
		argIdx++
	}
	if f.IsFeatured != nil {
		conditions = append(conditions, fmt.Sprintf("p.is_featured = $%d", argIdx))
		args = append(args, *f.IsFeatured)
		argIdx++
	}

	extraWhere := ""
	if len(conditions) > 0 {
		extraWhere = " AND " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products p WHERE p.is_active = true AND p.deleted_at IS NULL%s", extraWhere)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Paginate
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	offset := (f.Page - 1) * f.PageSize
	dataQuery := fmt.Sprintf(
		`SELECT p.id, p.name, p.slug, p.description, p.price, p.category_id, p.is_active, p.created_at, p.is_soon, p.is_featured
		 FROM products p
		 WHERE p.is_active = true AND p.deleted_at IS NULL%s
		 ORDER BY p.is_soon ASC, p.created_at DESC LIMIT $%d OFFSET $%d`,
		extraWhere, argIdx, argIdx+1,
	)
	args = append(args, f.PageSize, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Slug, &p.Description, &p.Price,
			&p.CategoryID, &p.IsActive, &p.CreatedAt, &p.IsSoon,
			&p.IsFeatured,
		); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}

	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *models.Product) error {
	query := `
		UPDATE products SET name=$1, slug=$2, description=$3, price=$4, stock=$5,
		       category_id=$6, image_url=$7, is_active=$8, updated_at=NOW()
		WHERE id=$9
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		p.Name, p.Slug, p.Description, p.Price, p.Stock,
		p.CategoryID, p.ImageURL, p.IsActive, p.ID,
	).Scan(&p.UpdatedAt)
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM products WHERE id = $1", id)
	return err
}

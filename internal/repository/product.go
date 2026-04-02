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
	           p.content, p.created_at, p.description, p.old_price
	    FROM products p
		left join brands b on b.id=p.brand_id
	    WHERE p.id = $1 AND p.is_active = true AND p.deleted_at IS NULL`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Sku, &p.Slug, &p.Price, &p.BrandId,&p.BrandName, &p.CategoryID,
		&p.Content, &p.CreatedAt, &p.Description, &p.OldPrice,
	)
	if err != nil {
		return nil, err
	}

	// Fetch all images for this product
	p.Images, err = r.getProductImages(ctx, []int64{id})
	if err != nil {
		p.Images = []string{}
	}

	return &p, nil
}

// getProductImages fetches image URLs for the given product IDs from the media table.
// Returns a slice of URLs for a single product, or is used internally to populate multiple products.
func (r *ProductRepository) getProductImages(ctx context.Context, productIDs []int64) ([]string, error) {
	if len(productIDs) == 0 {
		return []string{}, nil
	}

	query := `SELECT 'https://admin.origo.kz/storage/' || m.id || '/' || m.file_name
	    FROM media m
	    WHERE m.model_id = ANY($1)
	      AND m.model_type = 'App\Models\Product'
	    ORDER BY m.order_column`

	rows, err := r.db.Query(ctx, query, productIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		images = append(images, url)
	}
	if images == nil {
		images = []string{}
	}
	return images, nil
}

// getProductImagesMap fetches image URLs for multiple products and returns a map of productID -> []imageURL.
func (r *ProductRepository) getProductImagesMap(ctx context.Context, productIDs []int64) (map[int64][]string, error) {
	result := make(map[int64][]string)
	if len(productIDs) == 0 {
		return result, nil
	}

	query := `SELECT m.model_id, 'https://admin.origo.kz/storage/' || m.id || '/' || m.file_name
	    FROM media m
	    WHERE m.model_id = ANY($1)
	      AND m.model_type = 'App\Models\Product'
	    ORDER BY m.model_id, m.order_column`

	rows, err := r.db.Query(ctx, query, productIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var modelID int64
		var url string
		if err := rows.Scan(&modelID, &url); err != nil {
			return nil, err
		}
		result[modelID] = append(result[modelID], url)
	}
	return result, nil
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
		`SELECT p.id, p.name, p.slug, p.description, p.price, p.category_id, p.is_active, p.created_at
		 FROM products p
		 WHERE p.is_active = true AND p.deleted_at IS NULL%s
		 ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d`,
		extraWhere, argIdx, argIdx+1,
	)
	args = append(args, f.PageSize, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	var productIDs []int64
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Slug, &p.Description, &p.Price,
			&p.CategoryID, &p.IsActive, &p.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
		productIDs = append(productIDs, p.ID)
	}

	// Fetch images for all products in one query
	imagesMap, err := r.getProductImagesMap(ctx, productIDs)
	if err != nil {
		// If images fail, still return products with empty images
		for i := range products {
			products[i].Images = []string{}
		}
		return products, total, nil
	}

	for i := range products {
		if imgs, ok := imagesMap[products[i].ID]; ok {
			products[i].Images = imgs
		} else {
			products[i].Images = []string{}
		}
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

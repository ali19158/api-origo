package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/online-shop/internal/models"
)

type AttributeRepository struct {
	db *pgxpool.Pool
}

func NewAttributeRepository(db *pgxpool.Pool) *AttributeRepository {
	return &AttributeRepository{db: db}
}

// GetByProductID fetches all attributes with their values for a given product.
func (r *AttributeRepository) GetByProductID(ctx context.Context, productID int64) ([]models.ProductAttribute, error) {
	query := `SELECT a.id, a.name, a.slug, a.type, a.unit, a.sort_order, ap.value
	    FROM attribute_product ap
	    JOIN attributes a ON a.id = ap.attribute_id
	    WHERE ap.product_id = $1
	    ORDER BY a.sort_order, a.id`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attrs []models.ProductAttribute
	for rows.Next() {
		var pa models.ProductAttribute
		if err := rows.Scan(
			&pa.ID, &pa.Name, &pa.Slug, &pa.Type, &pa.Unit, &pa.SortOrder, &pa.Value,
		); err != nil {
			return nil, err
		}
		attrs = append(attrs, pa)
	}
	return attrs, nil
}

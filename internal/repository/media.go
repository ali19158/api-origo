package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/online-shop/internal/models"
)

type MediaRepository struct {
	db *pgxpool.Pool
}

func NewMediaRepository(db *pgxpool.Pool) *MediaRepository {
	return &MediaRepository{db: db}
}

// GetByModelIDs fetches media rows for the given model IDs and model type.
// modelType should be the Laravel model class, e.g. "App\Models\Product" or "App\Models\Category".
func (r *MediaRepository) GetByModelIDs(ctx context.Context, modelIDs []int64, modelType string) ([]models.MediaItem, error) {
	if len(modelIDs) == 0 {
		return nil, nil
	}

	query := `SELECT m.model_id, m.id, m.file_name
	    FROM media m
	    WHERE m.model_id = ANY($1)
	      AND m.model_type = $2
	    ORDER BY m.model_id, m.order_column`

	rows, err := r.db.Query(ctx, query, modelIDs, modelType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MediaItem
	for rows.Next() {
		var item models.MediaItem
		if err := rows.Scan(&item.ModelID, &item.MediaID, &item.FileName); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

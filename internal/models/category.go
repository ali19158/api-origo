package models

import (
	"encoding/json"
	"time"
)

type Category struct {
	ID            int64           `json:"id"`
	Name          json.RawMessage `json:"name"`
	Slug          string          `json:"slug"`
	Description   json.RawMessage `json:"description"`
	ParentID      *int64          `json:"parent_id,omitempty"`
	Preview       *string         `json:"preview,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	MediaID       *int64          `json:"-"`
	MediaFileName *string         `json:"-"`
}

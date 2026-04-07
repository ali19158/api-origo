package models

import "encoding/json"

// Attribute represents a product characteristic definition.
type Attribute struct {
	ID        int64           `json:"id"`
	Name      json.RawMessage `json:"name"`
	Slug      string          `json:"slug"`
	Type      string          `json:"type"`
	Unit      *string         `json:"unit"`
	SortOrder int             `json:"sort_order"`
}

// ProductAttribute represents a characteristic value attached to a product.
type ProductAttribute struct {
	Attribute
	Value string `json:"value"`
}

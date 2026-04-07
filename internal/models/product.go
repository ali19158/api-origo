package models

import (
	"encoding/json"
	"time"
)

type Product struct {
	ID          int64              `json:"id"`
	Name        json.RawMessage    `json:"name"`
	Slug        string             `json:"slug"`
	Description *json.RawMessage   `json:"description"`
	Price       *json.RawMessage   `json:"price"`
	Stock       int                `json:"stock"`
	CategoryID  int64              `json:"category_id"`
	ImageURL    string             `json:"image_url,omitempty"`
	IsActive    bool               `json:"is_active"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty"`
	BrandName   *json.RawMessage   `json:"brand_name"`
	IsSoon      bool               `json:"is_soon"`
	IsFeatured  bool               `json:"is_featured"`
	Sku         string             `json:"sku"`
	BrandId     int64              `json:"brand_id"`
	Content     *json.RawMessage   `json:"content,omitempty"`
	OldPrice    *json.RawMessage   `json:"old_price"`
	FileName    string             `json:"file_name"`
	Images      []string           `json:"images"`
	Attributes  []ProductAttribute `json:"attributes,omitempty"`
}

type ProductFilter struct {
	CategoryID *int64
	MinPrice   *float64
	MaxPrice   *float64
	Search     *string
	IsFeatured *bool
	Page       int
	PageSize   int
}

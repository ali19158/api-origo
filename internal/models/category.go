package models

import (
	"encoding/json"
	"time"
)

type Category struct {
	ID          int64           `json:"id"`
	Name        json.RawMessage `json:"name"`
	Slug        string          `json:"slug"`
	Description json.RawMessage `json:"description"`
	ParentID    *int64          `json:"parent_id,omitempty"`
	Preview     *string         `json:"preview,omitempty"`
	IsSoon      bool            `json:"is_soon"`
	CreatedAt   time.Time       `json:"created_at"`
}

// CategoryNode is used for the tree endpoint — same fields as Category plus nested children.
type CategoryNode struct {
	ID          int64           `json:"id"`
	Name        json.RawMessage `json:"name"`
	Slug        string          `json:"slug"`
	Description json.RawMessage `json:"description"`
	ParentID    *int64          `json:"parent_id,omitempty"`
	Preview     *string         `json:"preview,omitempty"`
	IsSoon      bool            `json:"is_soon"`
	CreatedAt   time.Time       `json:"created_at"`
	Children    []CategoryNode  `json:"children"`
}

// CategoryToNode converts a Category to a CategoryNode (without children).
func CategoryToNode(c Category) CategoryNode {
	return CategoryNode{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		ParentID:    c.ParentID,
		Preview:     c.Preview,
		IsSoon:      c.IsSoon,
		CreatedAt:   c.CreatedAt,
		Children:    []CategoryNode{},
	}
}

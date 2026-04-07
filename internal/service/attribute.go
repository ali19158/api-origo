package service

import (
	"context"

	"github.com/online-shop/internal/models"
	"github.com/online-shop/internal/repository"
)

type AttributeService struct {
	repo *repository.AttributeRepository
}

func NewAttributeService(repo *repository.AttributeRepository) *AttributeService {
	return &AttributeService{repo: repo}
}

// GetByProductID returns all characteristics for a given product.
func (s *AttributeService) GetByProductID(ctx context.Context, productID int64) ([]models.ProductAttribute, error) {
	attrs, err := s.repo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if attrs == nil {
		attrs = []models.ProductAttribute{}
	}
	return attrs, nil
}

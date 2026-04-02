package service

import (
	"context"

	"github.com/online-shop/internal/models"
	"github.com/online-shop/internal/repository"
)

type ProductService struct {
	repo     *repository.ProductRepository
	adminURL string
}

func NewProductService(repo *repository.ProductRepository, adminURL string) *ProductService {
	return &ProductService{repo: repo, adminURL: adminURL}
}

func (s *ProductService) Create(ctx context.Context, p *models.Product) error {
	return s.repo.Create(ctx, p)
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context, filter models.ProductFilter) ([]models.Product, int, error) {
	return s.repo.List(ctx, filter)
}

func (s *ProductService) Update(ctx context.Context, p *models.Product) error {
	return s.repo.Update(ctx, p)
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

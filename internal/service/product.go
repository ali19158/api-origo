package service

import (
	"context"

	"github.com/online-shop/internal/models"
	"github.com/online-shop/internal/repository"
)

const productModelType = `App\Models\Product`

type ProductService struct {
	repo     *repository.ProductRepository
	mediaSvc *MediaService
}

func NewProductService(repo *repository.ProductRepository, mediaSvc *MediaService) *ProductService {
	return &ProductService{repo: repo, mediaSvc: mediaSvc}
}

// enrichImages fetches media via MediaService and maps them to product image URLs.
func (s *ProductService) enrichImages(ctx context.Context, products []models.Product) {
	if len(products) == 0 {
		return
	}

	ids := make([]int64, len(products))
	for i, p := range products {
		ids[i] = p.ID
	}

	imagesMap, err := s.mediaSvc.GetImageURLs(ctx, ids, productModelType)
	if err != nil {
		for i := range products {
			products[i].Images = []string{}
		}
		return
	}

	for i := range products {
		if imgs, ok := imagesMap[products[i].ID]; ok {
			products[i].Images = imgs
		} else {
			products[i].Images = []string{}
		}
	}
}

func (s *ProductService) Create(ctx context.Context, p *models.Product) error {
	return s.repo.Create(ctx, p)
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	products := []models.Product{*p}
	s.enrichImages(ctx, products)
	*p = products[0]
	return p, nil
}

func (s *ProductService) List(ctx context.Context, filter models.ProductFilter) ([]models.Product, int, error) {
	products, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	s.enrichImages(ctx, products)
	return products, total, nil
}

func (s *ProductService) Update(ctx context.Context, p *models.Product) error {
	return s.repo.Update(ctx, p)
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

package service

import (
	"context"
	"fmt"

	"github.com/online-shop/internal/models"
	"github.com/online-shop/internal/repository"
)

type CategoryService struct {
	repo     *repository.CategoryRepository
	adminURL string
}

func NewCategoryService(repo *repository.CategoryRepository, adminURL string) *CategoryService {
	return &CategoryService{repo: repo, adminURL: adminURL}
}

func (s *CategoryService) enrichPreview(c *models.Category) {
	if c.MediaID != nil && c.MediaFileName != nil {
		url := fmt.Sprintf("%s/storage/%d/%s", s.adminURL, *c.MediaID, *c.MediaFileName)
		c.Preview = &url
	}
}

func (s *CategoryService) Create(ctx context.Context, c *models.Category) error {
	return s.repo.Create(ctx, c)
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (*models.Category, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.enrichPreview(c)
	return c, nil
}

func (s *CategoryService) List(ctx context.Context) ([]models.Category, error) {
	cats, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range cats {
		s.enrichPreview(&cats[i])
	}
	return cats, nil
}

func (s *CategoryService) Update(ctx context.Context, c *models.Category) error {
	return s.repo.Update(ctx, c)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

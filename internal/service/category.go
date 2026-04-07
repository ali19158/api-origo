package service

import (
	"context"

	"github.com/online-shop/internal/models"
	"github.com/online-shop/internal/repository"
)

const categoryModelType = `App\Models\Category`

type CategoryService struct {
	repo     *repository.CategoryRepository
	mediaSvc *MediaService
}

func NewCategoryService(repo *repository.CategoryRepository, mediaSvc *MediaService) *CategoryService {
	return &CategoryService{repo: repo, mediaSvc: mediaSvc}
}

// enrichPreviews fetches preview images for categories via MediaService.
func (s *CategoryService) enrichPreviews(ctx context.Context, categories []models.Category) {
	if len(categories) == 0 {
		return
	}

	ids := make([]int64, len(categories))
	for i, c := range categories {
		ids[i] = c.ID
	}

	previewMap, err := s.mediaSvc.GetFirstImageURL(ctx, ids, categoryModelType)
	if err != nil {
		return
	}

	for i := range categories {
		if url, ok := previewMap[categories[i].ID]; ok {
			categories[i].Preview = &url
		}
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
	cats := []models.Category{*c}
	s.enrichPreviews(ctx, cats)
	*c = cats[0]
	return c, nil
}

func (s *CategoryService) List(ctx context.Context) ([]models.Category, error) {
	cats, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	s.enrichPreviews(ctx, cats)
	return cats, nil
}

func (s *CategoryService) Update(ctx context.Context, c *models.Category) error {
	return s.repo.Update(ctx, c)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

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
	cats, err := s.repo.List(ctx, true)
	if err != nil {
		return nil, err
	}
	s.enrichPreviews(ctx, cats)
	return cats, nil
}

// ListRootCategories returns only root categories (parent_id IS NULL) with previews.
func (s *CategoryService) ListRootCategories(ctx context.Context) ([]models.Category, error) {
	cats, err := s.repo.ListRootCategories(ctx)
	if err != nil {
		return nil, err
	}
	s.enrichPreviews(ctx, cats)
	return cats, nil
}

// ListSubcategories returns direct children of the given parent category with previews.
func (s *CategoryService) ListSubcategories(ctx context.Context, parentID int64) ([]models.Category, error) {
	cats, err := s.repo.ListSubcategories(ctx, parentID)
	if err != nil {
		return nil, err
	}
	s.enrichPreviews(ctx, cats)
	return cats, nil
}

// ListTree returns the full category tree: root categories with nested children.
// Fetches all categories in one query and builds the tree in memory.
func (s *CategoryService) ListTree(ctx context.Context) ([]models.CategoryNode, error) {
	cats, err := s.repo.List(ctx, false)
	if err != nil {
		return nil, err
	}
	s.enrichPreviews(ctx, cats)

	// Index by ID for quick lookup
	nodeMap := make(map[int64]*models.CategoryNode, len(cats))
	for _, c := range cats {
		node := models.CategoryToNode(c)
		nodeMap[c.ID] = &node
	}

	// Build tree
	var roots []models.CategoryNode
	for _, c := range cats {
		if c.ParentID == nil {
			roots = append(roots, *nodeMap[c.ID])
		} else {
			if parent, ok := nodeMap[*c.ParentID]; ok {
				parent.Children = append(parent.Children, *nodeMap[c.ID])
			}
		}
	}

	// Update roots with their filled children from nodeMap
	for i, root := range roots {
		if updated, ok := nodeMap[root.ID]; ok {
			roots[i] = *updated
		}
	}

	if roots == nil {
		roots = []models.CategoryNode{}
	}

	return roots, nil
}

func (s *CategoryService) Update(ctx context.Context, c *models.Category) error {
	return s.repo.Update(ctx, c)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

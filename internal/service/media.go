package service

import (
	"context"
	"fmt"

	"github.com/online-shop/internal/models"
	"github.com/online-shop/internal/repository"
)

type MediaService struct {
	repo     *repository.MediaRepository
	adminURL string
}

func NewMediaService(repo *repository.MediaRepository, adminURL string) *MediaService {
	return &MediaService{repo: repo, adminURL: adminURL}
}

// BuildURL constructs a full URL for a media item.
func (s *MediaService) BuildURL(mediaID int64, fileName string) string {
	return fmt.Sprintf("%s/storage/%d/%s", s.adminURL, mediaID, fileName)
}

// GetImageURLs fetches media for the given model IDs and returns a map of modelID -> []imageURL.
func (s *MediaService) GetImageURLs(ctx context.Context, modelIDs []int64, modelType string) (map[int64][]string, error) {
	result := make(map[int64][]string)
	if len(modelIDs) == 0 {
		return result, nil
	}

	items, err := s.repo.GetByModelIDs(ctx, modelIDs, modelType)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		url := s.BuildURL(item.MediaID, item.FileName)
		result[item.ModelID] = append(result[item.ModelID], url)
	}
	return result, nil
}

// GetFirstImageURL fetches media for the given model IDs and returns a map of modelID -> first imageURL.
// Useful for entities that only need a single preview image (e.g. categories).
func (s *MediaService) GetFirstImageURL(ctx context.Context, modelIDs []int64, modelType string) (map[int64]string, error) {
	result := make(map[int64]string)
	if len(modelIDs) == 0 {
		return result, nil
	}

	items, err := s.repo.GetByModelIDs(ctx, modelIDs, modelType)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if _, exists := result[item.ModelID]; !exists {
			result[item.ModelID] = s.BuildURL(item.MediaID, item.FileName)
		}
	}
	return result, nil
}

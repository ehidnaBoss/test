package service

import (
	models "zero-agency-tambov/internal/entity"
	"zero-agency-tambov/internal/repository"
)

type NewsService interface {
	UpdateNews(news *models.News) error
	GetAllNews(limit, offset int) ([]*models.News, error)
	GetNewsWithCategories(id int64) (*models.News, error)
}

type newsService struct {
	repo repository.NewsRepository
}

func NewNewsService(repo repository.NewsRepository) NewsService {
	return &newsService{repo: repo}
}

func (s *newsService) UpdateNews(news *models.News) error {
	return s.repo.UpdateNews(news)
}

func (s *newsService) GetAllNews(limit, offset int) ([]*models.News, error) {
	return s.repo.GetAllNews(limit, offset)
}

func (s *newsService) GetNewsWithCategories(id int64) (*models.News, error) {
	return s.repo.GetNewsWithCategories(id)
}

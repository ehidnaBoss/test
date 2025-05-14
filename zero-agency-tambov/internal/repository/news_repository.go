package repository

import (
	"database/sql"
	models "zero-agency-tambov/internal/entity"
	"zero-agency-tambov/logger"

	"github.com/lib/pq"
)

type NewsRepository interface {
	UpdateNews(news *models.News) error
	GetAllNews(limit, offset int) ([]*models.News, error)
	GetNewsWithCategories(id int64) (*models.News, error)
}

type newsRepository struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) UpdateNews(news *models.News) error {
	logger.Log.Infof("Starting to update news with ID: %d", news.Id)
	query := "UPDATE news SET Title = COALESCE(NULLIF($1, ''), Title), Content = COALESCE(NULLIF($2, ''), Content) WHERE Id = $3"
	_, err := r.db.Exec(query, news.Title, news.Content, news.Id)
	if err != nil {
		logger.Log.Errorf("Error updating news with ID %d: %v", news.Id, err)
		return err
	}
	logger.Log.Infof("Successfully updated news with ID: %d", news.Id)

	if len(news.Categories) > 0 {
		logger.Log.Infof("Updating categories for news with ID: %d", news.Id)
		_, err := r.db.Exec("DELETE FROM NewsCategories WHERE news_id = $1", news.Id)
		if err != nil {
			logger.Log.Errorf("Error deleting categories for news with ID %d: %v", news.Id, err)
			return err
		}

		for _, categoryID := range news.Categories {
			_, err := r.db.Exec("INSERT INTO NewsCategories (news_id, category_id) VALUES ($1, $2)", news.Id, categoryID)
			if err != nil {
				logger.Log.Errorf("Error inserting category %d for news with ID %d: %v", categoryID, news.Id, err)
				return err
			}
			logger.Log.Infof("Category %d added for news with ID: %d", categoryID, news.Id)
		}
	}
	return nil
}

func (r *newsRepository) GetAllNews(limit, offset int) ([]*models.News, error) {
	logger.Log.Infof("Выборка всех новостей с ограничением %d и смещением %d", limit, offset)
	query := "SELECT Id, Title, Content FROM public.news LIMIT $1 OFFSET $2"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		logger.Log.Errorf("Ошибка выборки новосте: %v", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var newsList []*models.News
	for rows.Next() {
		var news models.News
		if err := rows.Scan(&news.Id, &news.Title, &news.Content); err != nil {
			logger.Log.Errorf("Ошибка при сканировании строки в поисках новостей: %v", err)
			return nil, err
		}
		newsList = append(newsList, &news)
	}
	logger.Log.Infof("Успешно получены % d новостных сообщений", len(newsList))
	return newsList, nil
}

func (r *newsRepository) GetNewsWithCategories(newsId int64) (*models.News, error) {
	query := `
		SELECT n.Id, n.Title, n.Content, 
			array_agg(c.CategoryId) AS Categories
		FROM public.news n
		LEFT JOIN public.newscategories c ON c.NewsId = n.Id
		WHERE n.Id = $1
		GROUP BY n.Id
	`
	row := r.db.QueryRow(query, newsId)

	var news models.News
	var categories []int64
	err := row.Scan(&news.Id, &news.Title, &news.Content, pq.Array(&categories))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	news.Categories = categories
	return &news, nil
}

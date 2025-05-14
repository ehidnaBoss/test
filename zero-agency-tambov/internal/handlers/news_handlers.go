package handlers

import (
	"strconv"

	models "zero-agency-tambov/internal/entity"
	"zero-agency-tambov/internal/service"
	"zero-agency-tambov/logger"

	"github.com/gofiber/fiber/v2"
)

type NewsHandler struct {
	services service.NewsService
}

func NewNewsHandler(services service.NewsService) *NewsHandler {
	return &NewsHandler{services: services}
}

func (h *NewsHandler) UpdateNews(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		logger.Log.Errorf("Invalid ID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var news models.News
	if err := c.BodyParser(&news); err != nil {
		logger.Log.Errorf("Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	news.Id = id

	logger.Log.Infof("Updating news with ID: %d", id)

	if err := h.services.UpdateNews(&news); err != nil {
		logger.Log.Errorf("Error updating news with ID %d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Log.Infof("Successfully updated news with ID: %d", id)
	return c.JSON(fiber.Map{"message": "News updated successfully"})
}

func (h *NewsHandler) GetAllNews(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	logger.Log.Infof("Fetching news with limit: %d and offset: %d", limit, offset)

	newsList, err := h.services.GetAllNews(limit, offset)
	if err != nil {
		logger.Log.Errorf("Error fetching news: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.Log.Infof("Successfully fetched %d news items", len(newsList))
	return c.JSON(newsList)
}

func (h *NewsHandler) GetNewsById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	news, err := h.services.GetNewsWithCategories(int64(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if news == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "News not found"})
	}

	return c.JSON(news)
}

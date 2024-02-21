package service

import (
	"context"

	"market/api/models"
	"market/pkg/logger"
	"market/storage"
)

type categoryService struct {
	storage storage.IStorage
	log     logger.ILogger
}

func NewCategoryService(storage storage.IStorage, log logger.ILogger) categoryService {
	return categoryService{
		storage: storage,
		log:     log,
	}
}

func (c categoryService) Create(ctx context.Context, createCategory models.CreateCategory) (models.Category, error) {
	c.log.Info("category create service layer", logger.Any("category", createCategory))

	pKey, err := c.storage.Category().Create(ctx, createCategory)
	if err != nil {
		c.log.Error("ERROR in service layer while creating category", logger.Error(err))
		return models.Category{}, err
	}

	category, err := c.storage.Category().GetByID(ctx, models.PrimaryKey{
		ID: pKey,
	})
	if err != nil {
		c.log.Error("ERROR in service layer while getting category", logger.Error(err))
		return models.Category{}, err
	}

	return category, nil
}
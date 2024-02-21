package service

import (
	"market/pkg/logger"
	"market/storage"
)

type IServiceManager interface {
	Basket() basketService
	Category() categoryService

}

type Service struct {
	basketService basketService
	categoryService categoryService
}

func New(storage storage.IStorage,  log logger.ILogger) Service {
	services := Service{}

	services.basketService = NewBasketService(storage, log)

	return  services
}

func (s Service) Basket() basketService {
	return s.basketService
}

func (s Service) Category() categoryService {
	return s.categoryService
}
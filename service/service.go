package service

import "market/storage"

type IServiceManager interface {
	Basket() basketService
}

type Service struct {
	basketService basketService
}

func New(storage storage.IStorage) Service {
	services := Service{}

	services.basketService = NewBasketService(storage)

	return  services
}

func (s Service) Basket() basketService {
	return s.basketService
}
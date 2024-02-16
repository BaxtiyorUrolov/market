package service

import "market/storage"

type saleService struct {
	storage storage.IStorage
}

func NewSaleService(storage storage.IStorage) saleService {
	return saleService{
		storage: storage,
	}
}
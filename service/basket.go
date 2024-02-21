package service

import (
	"context"
	"market/api/models"
	"market/pkg/logger"
	"market/storage"
)

type basketService struct {
	storage storage.IStorage
	log     logger.ILogger
}

func NewBasketService(storage storage.IStorage, log  logger.ILogger) basketService {
	return basketService{
		storage: storage,
		log: log,
	}
}

func (b basketService) Create(ctx context.Context, createBasket models.CreateBasket) (models.Basket, error) {

	count, err := b.storage.Repository().ProductByID(context.Background(), createBasket.ProductID)


	// Agar olinmoqchi bulgan product omborda bulmasa habar berish

	if count < createBasket.Quantity {
		b.log.Error("We don't have enough product", logger.Error(err))
		return models.Basket{}, err
	}

	product, err := b.storage.Product().GetByID(context.Background(), createBasket.ProductID)
	if err != nil {
		b.log.Error("Error in service layer while getting product ByID for Basket", logger.Error(err))
		return models.Basket{}, err
	}

	totalPrice := product.Price * createBasket.Quantity

	baskets, err := b.storage.Basket().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  100,
		Search: createBasket.SaleID,
	})
	if err != nil {
		b.log.Error("Error in service layer while getting baskets by SaleID for create Basket", logger.Error(err))
		return models.Basket{}, err
	}

	// Agar yaratilmoqchi bo'lgan basketdan bazada mavjud bo'lsa uning sonini o'zgartirish

	for _, basket := range baskets.Baskets {
		if basket.ProductID == createBasket.ProductID {
			if count < basket.Quantity + createBasket.Quantity {
				b.log.Error("We don't have enough product")
				return models.Basket{}, nil
			}
			updateBasket := models.UpdateBasket{
				ID:        basket.ID,
				SaleID:    basket.SaleID,
				ProductID: basket.ProductID,
				Quantity:  basket.Quantity + createBasket.Quantity,
				Price:     basket.Price + totalPrice,
			}
			if _, err := b.storage.Basket().Update(context.Background(), updateBasket); err != nil {
				b.log.Error("Error in service layer when adding baskets", logger.Error(err))
				return models.Basket{}, err
			}
			return models.Basket{}, err
		}
	}

	createBasket.Price = totalPrice

	id, err :=  b.storage.Basket().Create(context.Background(), createBasket)
	if err != nil {
		b.log.Error("Error in service layer when creating basket", logger.Error(err))
		return models.Basket{}, err
	}

	createdBasket, err := b.storage.Basket().GetByID(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		b.log.Error("Error in service layer when getting basket by id for create basket", logger.Error(err))
		return createdBasket, err
	}

	return createdBasket, nil

}

func (b basketService) Get(ctx context.Context, id string) (models.Basket, error) {
	basket, err := b.storage.Basket().GetByID(ctx, models.PrimaryKey{ID: id})
	if err != nil {
		b.log.Error("error in service layer while getting by id", logger.Error(err))
		return models.Basket{}, err
	}

	return basket, nil
}

func (b basketService) GetList(ctx context.Context, request models.GetListRequest) (models.BasketsResponse, error) {
	b.log.Info("basket get list service layer", logger.Any("basket", request))

	baskets, err := b.storage.Basket().GetList(ctx, request)
	if err != nil {
		b.log.Error("error in service layer  while getting list", logger.Error(err))
		return models.BasketsResponse{}, err
	}

	return baskets, nil
}

func (b basketService) Update(ctx context.Context, basket models.UpdateBasket) (models.Basket, error) {
	id, err := b.storage.Basket().Update(ctx, basket)
	if err != nil {
		b.log.Error("error in service layer while updating", logger.Error(err))
		return models.Basket{}, err
	}

	updatedBasket, err := b.storage.Basket().GetByID(context.Background(), models.PrimaryKey{ID: id})
	if err != nil {
		b.log.Error("error in service layer while getting basket by id", logger.Error(err))
		return models.Basket{}, err
	}

	return updatedBasket, nil
}

func (b basketService) Delete(ctx context.Context, key models.PrimaryKey) error {
	err := b.storage.Basket().Delete(ctx, key)

	return err
}

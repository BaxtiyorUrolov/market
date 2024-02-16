package service

import (
	"context"
	"fmt"
	"market/api/models"
	"market/storage"
)

type basketService struct {
	storage storage.IStorage
}

func NewBasketService(storage storage.IStorage) basketService {
	return basketService{
		storage: storage,
	}
}

func (b basketService) Create(ctx context.Context, createBasket models.CreateBasket) (models.Basket, error) {

	count, err := b.storage.Repository().ProductByID(context.Background(), createBasket.ProductID)


	// Agar olinmoqchi bulgan product omborda bulmasa habar berish

	if count < createBasket.Quantity {
		fmt.Println("We don't have enough product")
		return models.Basket{}, err
	}

	product, err := b.storage.Product().GetByID(context.Background(), createBasket.ProductID)
	if err != nil {
		fmt.Println("Error in service layer while getting product ByID for Basket", err.Error())
		return models.Basket{}, err
	}

	totalPrice := product.Price * createBasket.Quantity

	baskets, err := b.storage.Basket().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  100,
		Search: createBasket.SaleID,
	})
	if err != nil {
		fmt.Println("Error in service layer while getting baskets by SaleID for create Basket", err.Error())
		return models.Basket{}, err
	}

	// Agar yaratilmoqchi bo'lgan basketdan bazada mavjud bo'lsa uning sonini o'zgartirish

	for _, basket := range baskets.Baskets {
		if basket.ProductID == createBasket.ProductID {
			if count < basket.Quantity + createBasket.Quantity {
				fmt.Println("We don't have enough product")
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
				fmt.Println("Error in service layer when adding baskets", err.Error())
				return models.Basket{}, err
			}
			return models.Basket{}, err
		}
	}

	createBasket.Price = totalPrice

	id, err :=  b.storage.Basket().Create(context.Background(), createBasket)
	if err != nil {
		fmt.Println("Error in service layer when creating basket", err.Error())
		return models.Basket{}, err
	}

	createdBasket, err := b.storage.Basket().GetByID(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		fmt.Println("Error in service layer when getting basket by id for create basket")
		return createdBasket, err
	}

	return createdBasket, nil

}

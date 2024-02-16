package models


type Basket struct {
	ID         string     `json:"id"`
	SaleID     string     `json:"sale_id"`
	ProductID  string     `json:"product_id"`
	Quantity   int        `json:"quantity"`
	Price      int        `json:"price"`
	CreatedAt  string  	  `json:"created_at"`
	UpdatedAt  string	  `json:"updated_at"`
}

type CreateBasket struct {
	SaleID     string     `json:"sale_id"`
	ProductID  string     `json:"product_id"`
	Quantity   int        `json:"quantity"`
	Price      int        `json:"-"`
}

type UpdateBasket struct {
	ID         string     `json:"id"`
	SaleID     string     `json:"sale_id"`
	ProductID  string     `json:"product_id"`
	Quantity   int        `json:"quantity"`
	Price      int        `json:"price"`
}

type BasketsResponse struct {
	Baskets    []Basket   `json:"basket"`
	Count      int        `json:"count"`
}
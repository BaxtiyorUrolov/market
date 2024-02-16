package models

type Product struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Price      int       `json:"price"`
	Barcode    int       `json:"barcode"`
	CategoryID string    `json:"category_id"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string	 `json:"updated_at"`
}

type CreateProduct struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Barcode    int    `json:"barcode"`
	CategoryID string `json:"category_id"`
}

type UpdateProduct struct {
	ID         string    `json:"-"`
	Name       string    `json:"name"`
	Price      int       `json:"price"`
	CategoryID string    `json:"category_id"`
}

type ProductResponse struct {
	Products []Product
	Count    int
}

type ProductGetListRequest struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Name    string `json:"name"`
	Barcode int    `json:"barcode"`
}

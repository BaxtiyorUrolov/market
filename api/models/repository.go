package models

import "time"

type Repository struct {
	ID         string    `json:"id"`
	ProductID  string    `json:"product_id"`
	BranchID   string    `json:"branch_id"`
	Count      int       `json:"count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateRepository struct {
	ProductID string `json:"paroduct_id"`
	BranchID  string `json:"branch_id"`
	Count     int    `json:"count"`
}

type UpdateRepository struct {
	ID        string `json:"id"`
	ProductID string `json:"paroduct_id"`
	BranchID  string `json:"branch_id"`
	Count     int    `json:"count"`
}

type RepositoriesResponse struct {
	Repositories    []Repository `json:"repositories"`
	Count     int    `json:"count"`
}
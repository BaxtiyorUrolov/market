package postgres

import (
	"context"
	"fmt"
	"market/api/models"
	"market/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type saleRepo struct {
	db *pgxpool.Pool
}

func NewSaleRepo(db *pgxpool.Pool) storage.ISaleStorage {
	return saleRepo{db: db}
}

func (s saleRepo) Create(ctx context.Context, sale models.CreateSale) (string, error) {
	id := uuid.New()
	query := `insert into sales (id, branch_id, shop_assistant_id, cashier_id, client_name)
								values($1, $2, $3, $4, $5)`

	if _, err := s.db.Exec(ctx, query, id,
		sale.BranchID,
		sale.ShopAssistantID,
		sale.CashierID,
		sale.ClientName); err != nil {
		fmt.Println("error is while inserting data", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (s saleRepo) GetByID(ctx context.Context, id string) (models.Sale, error) {
	sale := models.Sale{}
	query := `select id, branch_id, shop_assistant_id, cashier_id, payment_type, price, status, client_name, 
					created_at, updated_at from sales where id = $1 and deleted_at is null`

	if err := s.db.QueryRow(ctx, query, id).Scan(
		&sale.ID,
		&sale.BranchID,
		&sale.ShopAssistantID,
		&sale.CashierID,
		&sale.PaymentType,
		&sale.Price,
		&sale.Status,
		&sale.ClientName,
		&sale.CreatedAt,
		&sale.UpdatedAt); err != nil {
		fmt.Println("error is while selecting by id", err.Error())
		return models.Sale{}, err
	}
	return sale, nil
}

func (s saleRepo) GetList(ctx context.Context, request models.GetListRequest) (models.SaleResponse, error) {
	var (
		page              = request.Page
		offset            = (page - 1) * request.Limit
		count             = 0
		query, countQuery string
		sales             = []models.Sale{}
		search            = request.Search
	)

	countQuery = `select count(1) from sales where deleted_at is null `
	if search != "" {
		countQuery += fmt.Sprintf(` AND client_name ilike '%%%s%%' `, search)
	}

	if err := s.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning count", err.Error())
		return models.SaleResponse{}, err
	}

	query = `select id, branch_id, shop_assistant_id, cashier_id, payment_type, price, status, client_name, 
					created_at, updated_at from sales where deleted_at is null `

	if search != "" {
		query += fmt.Sprintf(` AND client_name ilike '%%%s%%' `, search)
	}

	query += ` LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(ctx, query, request.Limit, offset)
	for rows.Next() {
		sale := models.Sale{}
		if err = rows.Scan(
			&sale.ID,
			&sale.BranchID,
			&sale.ShopAssistantID,
			&sale.CashierID,
			&sale.PaymentType,
			&sale.Price,
			&sale.Status,
			&sale.ClientName,
			&sale.CreatedAt,
			&sale.UpdatedAt); err != nil {
			fmt.Println("error is while scanning sales", err.Error())
			return models.SaleResponse{}, err
		}
		sales = append(sales, sale)
	}
	return models.SaleResponse{
		Sales: sales,
		Count: count,
	}, nil
}

func (s saleRepo) Update(ctx context.Context, sale models.UpdateSale) (string, error) {
	query := `UPDATE sales SET shop_assistant_id = $1, cashier_id = $2, payment_type = $3, 
				price = $4, status = $5, updated_at = NOW() WHERE id = $6`

	if _, err := s.db.Exec(ctx, query,
		sale.ShopAssistantID,
		sale.CashierID,
		sale.PaymentType,
		sale.Price,
		sale.Status,
		sale.ID,
	); err != nil {
		fmt.Println("error is while updating sale", err.Error())
		return "", err
	}
	return sale.ID, nil
}


func (s saleRepo) Delete(ctx context.Context, id string) error {
	query := `update sales set deleted_at = now() where id = $1`
	if _, err := s.db.Exec(ctx, query, id); err != nil {
		fmt.Println("error is while deleting sale", err.Error())
		return err
	}
	return nil
}

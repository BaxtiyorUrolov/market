package postgres

import (
	"context"
	"database/sql"
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
	query := `INSERT INTO sales (id, branch_id, shop_assistant_id, cashier_id, client_name)
								VALUES($1, $2, $3, $4, $5)`

	if _, err := s.db.Exec(ctx, query, id,
		sale.BranchID,
		sale.ShopAssistantID,
		sale.CashierID,
		sale.ClientName,
		); err != nil {
		fmt.Println("error is while inserting sale data", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (s saleRepo) GetByID(ctx context.Context, id string) (models.Sale, error) {
	var (
		updatedAt sql.NullTime
		paymentType sql.NullString
	)
	sale := models.Sale{}
	query := `SELECT id, branch_id, shop_assistant_id, cashier_id, payment_type, price, status, client_name, 
					created_at, updated_at FROM sales WHERE id = $1 and deleted_at = 0`

	if err := s.db.QueryRow(ctx, query, id).Scan(
		&sale.ID,
		&sale.BranchID,
		&sale.ShopAssistantID,
		&sale.CashierID,
		&paymentType,
		&sale.Price,
		&sale.Status,
		&sale.ClientName,
		&sale.CreatedAt,
		&updatedAt,
		); err != nil {
		fmt.Println("error is while selecting by id", err.Error())
		return models.Sale{}, err
	}

	if updatedAt.Valid {
		sale.UpdatedAt = updatedAt.Time
	}

	if paymentType.Valid {
		sale.PaymentType = paymentType.String
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
		updatedAt  		  sql.NullTime
	)

	countQuery = `SELECT COUNT(*) FROM sales WHERE deleted_at = 0 `
	if search != "" {
		countQuery += fmt.Sprintf(` AND client_name ILIKE '%%%s%%' `, search)
	}

	if err := s.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning count", err.Error())
		return models.SaleResponse{}, err
	}

	query = `SELECT id, branch_id, shop_assistant_id, cashier_id, payment_type, price, status, client_name, 
					created_at, updated_at FROM sales WHERE deleted_at = 0 `

	if search != "" {
		query += fmt.Sprintf(` AND client_name ilike '%%%s%%' `, search)
	}

	query += ` LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting all sales", err.Error())
		return models.SaleResponse{}, err
	}


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
			&updatedAt,
			); err != nil {
			fmt.Println("error is while scanning sales", err.Error())
			return models.SaleResponse{}, err
		}

		if updatedAt.Valid {
			sale.UpdatedAt = updatedAt.Time
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
	query := `UPDATE sales SET deleted_at = extract(epoch from current_timestamp) WHERE id = $1`
	if _, err := s.db.Exec(ctx, query, id); err != nil {
		fmt.Println("error is while deleting sale", err.Error())
		return err
	}
	return nil
}

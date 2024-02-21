package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"market/api/models"
	"market/pkg/logger"
	"market/storage"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepo struct {
	db *pgxpool.Pool
	log logger.ILogger
}

func NewProductRepo(db *pgxpool.Pool, log logger.ILogger) storage.IProducts {
	return productRepo{
		db: db,
		log: log,
	}
}

func (p productRepo) Create(ctx context.Context, product models.CreateProduct) (string, error) {
	id := uuid.New()
	query := `INSERT INTO products (id, name, price, barcode, category_id) VALUES($1, $2, $3, $4, $5)`
	if _, err := p.db.Exec(ctx, query,
		id,
		product.Name,
		product.Price,
		product.Barcode,
		product.CategoryID,
	); err != nil {
		fmt.Println("error is while inserting data", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (p productRepo) GetByID(ctx context.Context, id string) (models.Product, error) {
	var updatedAt, createdAt sql.NullString
	product := models.Product{}
	query := `SELECT id, name, price, barcode, category_id, created_at, updated_at 
	FROM products WHERE id = $1 AND deleted_at = 0`
	if err := p.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Barcode,
		&product.CategoryID,
		&createdAt,
		&updatedAt,
	); err != nil {
		fmt.Println("error is while scanning: ", err.Error())
		return models.Product{}, err
	}

	if createdAt.Valid {
		product.CreatedAt = createdAt.String
	}

	if updatedAt.Valid {
		product.UpdatedAt = updatedAt.String
	}

	return product, nil
}

func (p productRepo) GetList(ctx context.Context, request models.ProductGetListRequest) (models.ProductResponse, error) {
	var (
		page              = request.Page
		offset            = (page - 1) * request.Limit
		query, countQuery string
		count             = 0
		products          = []models.Product{}
		name              = request.Name
		barcode           = request.Barcode
		createdAt         sql.NullString
		updatedAt         sql.NullString
	)
	countQuery = `SELECT count(1) FROM products WHERE deleted_at = 0 `

	if name != "" {
		countQuery += fmt.Sprintf(` AND name ilike '%%%s%%' `, name)
	}

	if barcode != 0 {
		countQuery += ` AND barcode = ` + strconv.Itoa(barcode)
	}

	if err := p.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning count ....", err.Error())
		return models.ProductResponse{}, err
	}

	query = `SELECT  id, name, price, barcode, category_id, created_at, updated_at 
							FROM products WHERE deleted_at = 0 `

	if name != "" {
		query += fmt.Sprintf(` AND name ilike '%%%s%%' `, name)
	}
	if barcode != 0 {
		query += ` AND barcode = ` + strconv.Itoa(barcode)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := p.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting all products", err.Error())
		return models.ProductResponse{}, err
	}

	for rows.Next() {
		product := models.Product{}
		if err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Barcode,
			&product.CategoryID,
			&createdAt,
			&updatedAt,
		); err != nil {
			fmt.Println("error is while scanning category", err.Error())
			return models.ProductResponse{}, err
		}

		if createdAt.Valid {
			product.CreatedAt = createdAt.String
		}

		if updatedAt.Valid {
			product.UpdatedAt = updatedAt.String
		}

		products = append(products, product)
	}
	return models.ProductResponse{
		Products: products,
		Count:    count,
	}, nil

}

func (p productRepo) Update(ctx context.Context, product models.UpdateProduct) (string, error) {
	query := `UPDATE products SET name = $1, price = $2, category_id = $3, updated_at = now() 
									WHERE id = $4 AND deleted_at = 0`
	if _, err := p.db.Exec(ctx, query,
		&product.Name,
		&product.Price,
		&product.CategoryID,
		&product.ID); err != nil {
		fmt.Println("error is while updating", err.Error())
		return "", err
	}
	return product.ID, nil
}

func (p productRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE products SET deleted_at = extract(epoch from current_timestamp) WHERE id = $1`
	if _, err := p.db.Exec(ctx, query, &id); err != nil {
		fmt.Println("error is while deleting", err.Error())
		return err
	}
	return nil
}

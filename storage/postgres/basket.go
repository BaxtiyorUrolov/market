package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"market/api/models"
	"market/pkg/logger"
	"market/storage"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type basketRepo struct {
	DB *pgxpool.Pool
	log logger.ILogger
}

func NewBasketRepo(DB *pgxpool.Pool, log logger.ILogger) storage.IBasketRepo {
	return &basketRepo{
		DB: DB,
		log: log,
	}
}

func (s *basketRepo) Create(ctx context.Context, basket models.CreateBasket) (string, error) {
	id := uuid.New().String()
	createdAT := time.Now()

	if _, err := s.DB.Exec(ctx, `INSERT INTO baskets 
		(id, sale_id, product_id, quantity, price, created_at)
			VALUES($1, $2, $3, $4, $5, $6)`,
		id,
		basket.SaleID,
		basket.ProductID,
		basket.Quantity,
		basket.Price,
		createdAT,
	); err != nil {
		s.log.Error("Error while inserting data:", logger.Error(err))
		return "", err
	}

	return id, nil
}

func (s *basketRepo) GetByID(ctx context.Context, id models.PrimaryKey) (models.Basket, error) {
	var updatedAt, createdAt sql.NullString
	basket := models.Basket{}
	query := `SELECT id, sale_id, product_id, quantity, price, created_at, updated_at
				FROM baskets WHERE id = $1 AND  deleted_at = 0`
	err := s.DB.QueryRow(ctx, query, id.ID).Scan(
		&basket.ID,
		&basket.SaleID,
		&basket.ProductID,
		&basket.Quantity,
		&basket.Price,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		s.log.Error("Error while selecting basket by ID:", logger.Error(err))
		return models.Basket{}, err
	}

	if createdAt.Valid {
		basket.CreatedAt = createdAt.String
	}

	if updatedAt.Valid {
		basket.UpdatedAt = updatedAt.String
	}

	return basket, nil
}

func (s *basketRepo) GetList(ctx context.Context, request models.GetListRequest) (models.BasketsResponse, error) {
	var (
		page              = request.Page
		offset            = (page - 1) * request.Limit
		query, countQuery string
		baskets = []models.Basket{}
		count   int
		updatedAt, createdAt  sql.NullString
	)

	countQuery = `SELECT COUNT(*) FROM baskets WHERE deleted_at = 0`
	if request.Search != "" {
		countQuery += fmt.Sprintf(`  AND sale_id = '%s'`, request.Search)
	}

	err := s.DB.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		s.log.Error("Error while scanning count of baskets:", logger.Error(err))
		return models.BasketsResponse{}, err
	}

	query = `SELECT id, sale_id, product_id, quantity, price, created_at, updated_at
						FROM baskets WHERE deleted_at = 0`
	if request.Search != "" {
		query += fmt.Sprintf(` AND sale_id = '%s'`, request.Search)
	}
	query += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := s.DB.Query(ctx, query, request.Limit, offset)
	if err != nil {
		s.log.Error("Error while querying baskets:", logger.Error(err))
		return models.BasketsResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		basket := models.Basket{}
		err := rows.Scan(
			&basket.ID,
			&basket.SaleID,
			&basket.ProductID,
			&basket.Quantity,
			&basket.Price,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			s.log.Error("Error while scanning row of baskets:", logger.Error(err))
			return models.BasketsResponse{}, err
		}

		if createdAt.Valid {
			basket.CreatedAt = createdAt.String
		}

		if updatedAt.Valid {
			basket.UpdatedAt = updatedAt.String
		}

		baskets = append(baskets, basket)
	}

	return models.BasketsResponse{
		Baskets: baskets,
		Count:   count,
	}, nil
}

func (s *basketRepo) Update(ctx context.Context, basket models.UpdateBasket) (string, error) {
	query := `UPDATE baskets SET sale_id = $1, product_id = $2, quantity = $3, price = $4, updated_at = NOW() WHERE id = $5 AND deleted_at = 0`

	_, err := s.DB.Exec(ctx, query,
		&basket.SaleID,
		&basket.ProductID,
		&basket.Quantity,
		&basket.Price,
		&basket.ID,
	)
	if err != nil {
		s.log.Error("Error while updating Basket :", logger.Error(err))
		return "", err
	}

	return basket.ID, nil
}

func (b *basketRepo) Delete(ctx context.Context, key models.PrimaryKey) error {
	query := `update baskets set deleted_at = extract(epoch from current_timestamp) where id = $1`
	if rowsAffected, err := b.DB.Exec(ctx, query, key.ID); err != nil {
		if r := rowsAffected.RowsAffected(); r == 0 {
			b.log.Error("error is while deleting basket", logger.Error(err))
			return err
		}
		return err
	}
	return nil
}

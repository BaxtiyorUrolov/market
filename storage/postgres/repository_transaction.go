package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"market/api/models"
	"market/storage"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryTransactionRepo struct {
	DB *pgxpool.Pool
}

func NewRepositoryTransactionRepo(DB *pgxpool.Pool) storage.IRepositoryTransactionRepo {
	return &repositoryTransactionRepo{
		DB: DB,
	}
}

func (s *repositoryTransactionRepo) Create(ctx context.Context, rtransaction models.CreateRepositoryTransaction) (string, error) {
	id := uuid.New().String()
	createdAt := time.Now()

	if _, err := s.DB.Exec(ctx, `INSERT INTO repository_transactions
		(id, staff_id, product_id, repository_transaction_type, price, quantity, created_at)
			VALUES($1, $2, $3, $4, $5, $6, $7)`,
			id,
			rtransaction.StaffID,
			rtransaction.ProductID,
			rtransaction.RepositoryTransactionType,
			rtransaction.Price,
			rtransaction.Quantity,
			createdAt,
	); err != nil {
        log.Println("Error while inserting data:", err)
        return "", err
    }

    return id, nil
}

func (s *repositoryTransactionRepo) GetByID(ctx context.Context, id models.PrimaryKey) (models.RepositoryTransaction, error) {
	var updatedAt sql.NullTime
	rtransaction := models.RepositoryTransaction{}
	query := `SELECT id, staff_id, product_id, repository_transaction_type, price, quantity, created_at, updated_at FROM repository_transactions WHERE id = $1`

	err := s.DB.QueryRow(ctx, query, id.ID).Scan(
		&rtransaction.ID,
		&rtransaction.StaffID,
		&rtransaction.ProductID,
		&rtransaction.RepositoryTransactionType,
		&rtransaction.Price,
		&rtransaction.Quantity,
		&updatedAt,
	)
	if err != nil {
		log.Println("Error while selecting repository by ID:", err)
		return models.RepositoryTransaction{}, err
	}

	if updatedAt.Valid {
		rtransaction.UpdatedAt = updatedAt.Time
	}

	return rtransaction, nil
}

func (s *repositoryTransactionRepo) GetList(ctx context.Context, request models.GetListRequest) (models.RepositoryTransactionsResponse, error) {
	var (
		page              = request.Page
		offset            = (page - 1) * request.Limit
		query, countQuery string
		rtransactions = []models.RepositoryTransaction{}
		count  int
		updatedAt   sql.NullTime
	)

	countQuery = `SELECT COUNT(*) FROM repository_transactions where deleted_at = 0 `
	if request.Search != "" {
		countQuery += fmt.Sprintf(`WHERE quantity ILIKE '%%%s%%' or price ilike '%%%s%%'`, request.Search, request.Search)
	}

	err := s.DB.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		log.Println("Error while scanning count of repository_transactions:", err)
		return models.RepositoryTransactionsResponse{}, err
	}

	query = `SELECT id, staff_id, product_id, repository_transaction_type, price, quantity, created_at, updated_at FROM repository_transactions WHERE deleted_at = 0 `
	if request.Search != "" {
		query += fmt.Sprintf(` WHERE quantity ILIKE '%%%s%%' or price ilike '%%%s%%'`, request.Search, request.Search)
	}
	query += ` LIMIT $1 OFFSET $2`

	rows, err := s.DB.Query(ctx, query, request.Limit, offset)
	if err != nil {
		log.Println("Error while querying repository_transactions:", err)
		return models.RepositoryTransactionsResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		rtransaction := models.RepositoryTransaction{}
		err := rows.Scan(
			&rtransaction.ID,
			&rtransaction.StaffID,
			&rtransaction.ProductID,
			&rtransaction.RepositoryTransactionType,
			&rtransaction.Price,
			&rtransaction.Quantity,
			&updatedAt,
		)
		if err != nil {
			log.Println("Error while scanning row of repository_transactions:", err)
			return models.RepositoryTransactionsResponse{}, err
		}

		if updatedAt.Valid {
			rtransaction.UpdatedAt = updatedAt.Time
		}

		rtransactions = append(rtransactions, rtransaction)
	}

	return models.RepositoryTransactionsResponse{
		RepositoryTransactions: rtransactions,
		Count:  count,
	}, nil
}

func (s *repositoryTransactionRepo) Update(ctx context.Context, rtransaction models.UpdateRepositoryTransaction) (string, error) {
	query := `UPDATE repository_transactions SET staff_id = $1, product_id = $2, repository_transaction_type = $3, price = $4, quantity = $5, updated_at = NOW() WHERE id = $6`

	_, err := s.DB.Exec(ctx, query,
		&rtransaction.StaffID,
		&rtransaction.ProductID,
		&rtransaction.RepositoryTransactionType,
		&rtransaction.Price,
		&rtransaction.Quantity,
		&rtransaction.ID,
	)
	if err != nil {
		log.Println("Error while repository_transactions Repository :", err)
		return "", err
	}

	return rtransaction.ID, nil
}

func (s *repositoryTransactionRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE repository_transactions SET deleted_at = extract(epoch from current_timestamp) WHERE id = $1`

	_, err := s.DB.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error while deleting repository_transactions ", err)
		return err
	}

	return nil
}
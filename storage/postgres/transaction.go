package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"market/api/models"
	"market/storage"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) storage.ITransactionStorage {
	return transactionRepo{db: db}
}

func (t transactionRepo) Create(ctx context.Context, trans models.CreateTransaction) (string, error) {
	id := uuid.New()
	query := `INSERT INTO transactions 
    					(id, sale_id, staff_id, transaction_type, source_type, amount, description) 
						VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := t.db.Exec(ctx, query, id,
		trans.SaleID,
		trans.StaffID,
		trans.TransactionType,
		trans.SourceType,
		trans.Amount,
		trans.Description,
		); err != nil {
		fmt.Println("error is while inserting data", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (t transactionRepo) GetByID(ctx context.Context, id string) (models.Transaction, error) {
	var updatedAt sql.NullTime
	trans := models.Transaction{}
	query := `SELECT id, sale_id, staff_id, transaction_type, source_type, amount,
       						description, created_at, updated_at
							FROM transactions where deleted_at = 0 AND id = $1`
	if err := t.db.QueryRow(ctx, query, id).Scan(
		&trans.ID,
		&trans.SaleID,
		&trans.StaffID,
		&trans.TransactionType,
		&trans.SourceType,
		&trans.Amount,
		&trans.Description,
		&trans.CreatedAt,
		&updatedAt,
		); err != nil {
		fmt.Println("error is while selecting by id", err.Error())
		return models.Transaction{}, err
	}
	return trans, nil
}

func (t transactionRepo) GetList(ctx context.Context, request models.TransactionGetListRequest) (models.TransactionResponse, error) {
	var (
		page              = request.Page
		offset            = (page - 1) * request.Limit
		transactions      = []models.Transaction{}
		fromAmount        = request.FromAmount
		toAmount          = request.ToAmount
		count             = 0
		query, countQuery string
		updatedAt         sql.NullTime
	)

	countQuery = `SELECT COUNT(1) FROM transactions WHERE deleted_at = 0 `
	if fromAmount != 0 && toAmount != 0 {
		countQuery += fmt.Sprintf(` AND amount between %f and %f`, fromAmount, toAmount)
	} else if fromAmount != 0 {
		countQuery += ` AND amount <= ` + strconv.FormatFloat(fromAmount, 'f', 2, 64)
	} else {
		countQuery += ` AND amount >= ` + strconv.FormatFloat(toAmount, 'f', 2, 64)

	}
	if err := t.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning row", err.Error())
		return models.TransactionResponse{}, err
	}

	query = `SELECT id, sale_id, staff_id, transaction_type, source_type, amount,
       						description, created_at, updated_at FROM transactions WHERE deleted_at = 0 `

	if fromAmount != 0 && toAmount != 0 {
		query += fmt.Sprintf(` AND amount between %f and %f  order by amount asc, `, fromAmount, toAmount)
	} else if fromAmount != 0 {
		query += ` AND amount <= ` + strconv.FormatFloat(fromAmount, 'f', 2, 64) + `  order by amount asc, `
	} else {
		query += ` AND amount >= ` + strconv.FormatFloat(toAmount, 'f', 2, 64) + ` order by amount asc, `

	}

	query += ` created_at desc LIMIT $1 OFFSET $2 `

	rows, err := t.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting all from transactions", err.Error())
		return models.TransactionResponse{}, err
	}

	for rows.Next() {
		trans := models.Transaction{}
		if err = rows.Scan(
			&trans.ID,
			&trans.SaleID,
			&trans.StaffID,
			&trans.TransactionType,
			&trans.SourceType,
			&trans.Amount,
			&trans.Description,
			&trans.CreatedAt,
			&updatedAt,
			); err != nil {
			fmt.Println("error is while scanning rows", err.Error())
			return models.TransactionResponse{}, err
		}

		if updatedAt.Valid {
			trans.UpdatedAt = updatedAt.Time
		}

		transactions = append(transactions, trans)
	}
	return models.TransactionResponse{
		Transactions: transactions,
		Count:        count,
	}, nil
}

func (t transactionRepo) Update(ctx context.Context, transaction models.UpdateTransaction) (string, error) {
	query := `UPDATE transactions SET sale_id = $1, staff_id = $2, transaction_type = $3, source_type = $4, amount = $5,
								description = $6, updated_at = NOW() 
                    			WHERE id = $7 AND deleted_at = 0`
	if _, err := t.db.Exec(ctx, query,
		&transaction.SaleID,
		&transaction.StaffID,
		&transaction.TransactionType,
		&transaction.SourceType,
		&transaction.Amount,
		&transaction.Description,
		&transaction.ID,
		); err != nil {
		fmt.Println("error is while updating transaction: ", err.Error())
		return "", err
	}
	return transaction.ID, nil
}

func (t transactionRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE transactions SET deleted_at = extract(epoch from current_timestamp) WHERE id = $1`
	if _, err := t.db.Exec(ctx, query, id); err != nil {
		fmt.Println("error is while deleting transaction: ", err.Error())
		return err
	}
	return nil
}
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"market/api/models"
	"market/pkg/logger"
	"market/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type branchRepo struct {
	db *pgxpool.Pool
	log logger.ILogger
}

func NewBranchRepo(db *pgxpool.Pool, log logger.ILogger) storage.IBranchStorage {
	return branchRepo{
		db: db,
		log: log,
	}
}
func (b branchRepo) Create(ctx context.Context, branch models.CreateBranch) (string, error) {
	id := uuid.New()

	query := `INSERT INTO branches (id, name, address) 
									VALUES($1, $2, $3)`

	if _, err := b.db.Exec(ctx, query,
		id,
		branch.Name,
		branch.Address); err != nil {
		b.log.Error("error is while inserting data", logger.Error(err))
		return "", err
	}

	return id.String(), nil
}

func (b branchRepo) GetByID(ctx context.Context, id string) (models.Branch, error) {
	var updatedAt sql.NullTime
	branch := models.Branch{}
	query := `SELECT id, name, address, created_at, updated_at FROM branches WHERE id = $1 AND deleted_at = 0`
	if err := b.db.QueryRow(ctx, query, id).Scan(
		&branch.ID,
		&branch.Name,
		&branch.Address,
		&branch.CreatedAt,
		&updatedAt,
		); err != nil {
		b.log.Error("error is while selecting by id", logger.Error(err))
		return models.Branch{}, err
	}

	if updatedAt.Valid {
		branch.UpdatedAt = updatedAt.Time
	}

	return branch, nil
}

func (b branchRepo) GetList(ctx context.Context, request models.GetListRequest) (models.BranchResponse, error) {
	var (
		count             = 0
		branches          = []models.Branch{}
		query, countQuery string
		page              = request.Page
		offset            = (page - 1) * request.Limit
		search            = request.Search
		updatedAt         sql.NullTime
	)

	countQuery = `SELECT COUNT(1) FROM branches WHERE deleted_at = 0 `

	if search != "" {
		countQuery += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}

	if err := b.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		b.log.Error("error is while scanning count", logger.Error(err))
		return models.BranchResponse{}, err
	}

	query = `SELECT id, name, address, created_at, updated_at FROM branches WHERE deleted_at = 0 `
	if search != "" {
		query += fmt.Sprintf(` AND name ILIKE '%%%s%%' `, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := b.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		b.log.Error("error is while selecting * from branches", logger.Error(err))
		return models.BranchResponse{}, err
	}

	for rows.Next() {
		branch := models.Branch{}
		if err := rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.Address,
			&branch.CreatedAt,
			&updatedAt,
			); err != nil {
			b.log.Error("error is while scanning branch", logger.Error(err))
			return models.BranchResponse{}, err
		}

		if updatedAt.Valid {
			branch.UpdatedAt = updatedAt.Time
		}

		branches = append(branches, branch)
	}

	return models.BranchResponse{
		Branches: branches,
		Count:    count,
	}, err
}
func (b branchRepo) Update(ctx context.Context, branch models.UpdateBranch) (string, error) {
	query := `UPDATE branches SET name = $1, address = $2, updated_at = Now() WHERE id = $3 AND deleted_at = 0`

	if _, err := b.db.Exec(ctx, query,
		&branch.Name,
		&branch.Address,
		&branch.ID); err != nil {
		b.log.Error("error is while updating branch", logger.Error(err))
		return "", err
	}

	return branch.ID, nil
}
func (b branchRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE branches SET deleted_at = extract(epoch from current_timestamp) WHERE id = $1`

	if _, err := b.db.Exec(ctx, query, id); err != nil {
		b.log.Error("error is while deleting branches", logger.Error(err))
		return err
	}
	return nil
}

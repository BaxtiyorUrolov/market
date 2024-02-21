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

type categoryRepo struct {
	db *pgxpool.Pool
	log logger.ILogger
}

func NewCategoryRepo(db *pgxpool.Pool, log logger.ILogger) storage.ICategory {
	return categoryRepo{
		db: db,
		log: log,
	}
}

func (c categoryRepo) Create(ctx context.Context, category models.CreateCategory) (string, error) {
	id := uuid.New()
	query := `insert into categories (id, name, parent_id) values($1, $2, $3)`
	if _, err := c.db.Exec(ctx, query, id, category.Name, category.ParentID); err != nil {
		c.log.Error("error is while inserting data", logger.Error(err))
		return "", err
	}
	return id.String(), nil
}

func (c categoryRepo) GetByID(ctx context.Context, id models.PrimaryKey) (models.Category, error) {
	var updatedAt sql.NullTime
	category := models.Category{}
	query := `select id, name, parent_id, created_at, updated_at FROM categories WHERE id = $1 and deleted_at = 0`
	if err := c.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.ParentID,
		&category.CreatedAt,
		&updatedAt,
		); err != nil {
		c.log.Error("error is while selecting by id", logger.Error(err))
		return models.Category{}, err
	}

	if updatedAt.Valid {
		category.UpdatedAt = updatedAt.Time
	}

	return category, nil
}

func (c categoryRepo) GetList(ctx context.Context, request models.GetListRequest) (models.CategoryResponse, error) {
	var (
		page              = request.Page
		offset            = (page - 1) * request.Limit
		query, countQuery string
		count             = 0
		categories        = []models.Category{}
		search            = request.Search
		updatedAt		  sql.NullTime
	)
	countQuery = `SELECT count(1) FROM categories WHERE deleted_at = 0 `
	if search != "" {
		countQuery += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}
	if err := c.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		c.log.Error("error is while scanning count", logger.Error(err))
		return models.CategoryResponse{}, err
	}

	query = `SELECT id, name, parent_id, created_at, updated_at FROM categories WHERE deleted_at = 0 `
	if search != "" {
		query += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}

	query += ` LIMIT $1 OFFSET $2`
	rows, err := c.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		c.log.Error("error is while selecting all", logger.Error(err))
		return models.CategoryResponse{}, err
	}

	for rows.Next() {
		category := models.Category{}
		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.ParentID,
			&category.CreatedAt,
			&updatedAt,
			); err != nil {
			c.log.Error("error is while scanning category", logger.Error(err))
			return models.CategoryResponse{}, err
		}

		if updatedAt.Valid {
			category.UpdatedAt = updatedAt.Time
		}

		categories = append(categories, category)
	}
	return models.CategoryResponse{
		Categories: categories,
		Count:      count,
	}, nil
}

func (c categoryRepo) Update(ctx context.Context, category models.UpdateCategory) (string, error) {
	query := `UPDATE categories SET name = $1, parent_id = $2, updated_at = now() WHERE id = $3 AND deleted_at = 0`
	if _, err := c.db.Exec(ctx, query, &category.Name, &category.ParentID, &category.ID); err != nil {
		c.log.Error("error is while updating", logger.Error(err))
		return "", err
	}
	return category.ID, nil
}

func (c categoryRepo) Delete(ctx context.Context, id string) error {
	query := `update categories set deleted_at = extract(epoch FROM current_timestamp) WHERE id = $1`
	if _, err := c.db.Exec(ctx, query, id); err != nil {
		c.log.Error("error is while deleting", logger.Error(err))
		return err
	}
	return nil
}

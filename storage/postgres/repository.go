package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"market/api/models"
	"market/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryRepo struct {
	DB *pgxpool.Pool
}

func NewRepositoryRepo(DB *pgxpool.Pool) storage.IRepositoryRepo {
	return &repositoryRepo{
		DB: DB,
	}
}

func (s *repositoryRepo) Create(ctx context.Context, repository models.CreateRepository) (string, error) {
	id := uuid.New()

	if _, err := s.DB.Exec(ctx, `INSERT INTO repositories 
    (id, product_id, branch_id, count) 
        VALUES ($1, $2, $3, $4)`,
		id,
		repository.ProductID,
		repository.BranchID,
		repository.Count,
	); err != nil {
		log.Println("Error while inserting data:", err)
		return "", err
	}

	return id.String(), nil
}

func (s *repositoryRepo) GetByID(ctx context.Context, id models.PrimaryKey) (models.Repository, error) {
	var updatedAt sql.NullTime
	repository := models.Repository{}
	query := `SELECT id, product_id, branch_id, count, created_at, updated_at FROM repositories WHERE id = $1`
	err := s.DB.QueryRow(ctx, query, id.ID).Scan(
		&repository.ID,
		&repository.ProductID,
		&repository.BranchID,
		&repository.Count,
		&repository.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		log.Println("Error while selecting repository by ID:", err)
		return models.Repository{}, err
	}

	if updatedAt.Valid {
		repository.UpdatedAt = updatedAt.Time
	}

	return repository, nil
}

func (s *repositoryRepo) ProductByID(ctx context.Context, id string) (int, error) {
	var count int

	query := `SELECT count FROM repositories WHERE product_id = $1`
	err := s.DB.QueryRow(ctx, query, id).Scan(
		&count,
	)
	if err != nil {
		fmt.Println("Error while selecting product count: ", err)
		return count, err
	}

	return count, nil
}

func (s *repositoryRepo) GetList(ctx context.Context, request models.GetListRequest) (models.RepositoriesResponse, error) {
	var (
		page              = request.Page
		offset            = (page - 1) * request.Limit
		repositories = []models.Repository{}
		count        int
		updatedAt  		  sql.NullTime
	)

	countQuery := `SELECT COUNT(*) FROM repositories WHERE deleted_at IS NULL`
	if request.Search != "" {
		countQuery += fmt.Sprintf(` AND product_id = '%s'`, request.Search)
	}

	err := s.DB.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		log.Println("Error while scanning count of repositories:", err)
		return models.RepositoriesResponse{}, err
	}

	query := `SELECT id, product_id, branch_id, count, created_at, updated_at 
			  FROM repositories WHERE deleted_at IS NULL`
	if request.Search != "" {
		query += fmt.Sprintf(` AND product_id = '%s'`, request.Search)
	}
	query += ` LIMIT $1 OFFSET $2`

	rows, err := s.DB.Query(ctx, query, request.Limit, offset)
	if err != nil {
		log.Println("Error while querying repositories:", err)
		return models.RepositoriesResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		repository := models.Repository{}
		err := rows.Scan(
			&repository.ID,
			&repository.ProductID,
			&repository.BranchID,
			&repository.Count,
			&repository.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			log.Println("Error while scanning row of repositories:", err)
			return models.RepositoriesResponse{}, err
		}

		if updatedAt.Valid {
			repository.UpdatedAt = updatedAt.Time
		}

		repositories = append(repositories, repository)
	}

	return models.RepositoriesResponse{
		Repositories: repositories,
		Count:        count,
	}, nil
}

func (s *repositoryRepo) Update(ctx context.Context, repository models.UpdateRepository) (string, error) {
	query := `UPDATE repositories SET branch_id = $1, product_id = $2, count = $3, updated_at = NOW() WHERE id = $4`

	_, err := s.DB.Exec(ctx, query,
		repository.BranchID,
		repository.ProductID,
		repository.Count,
		repository.ID,
	)
	if err != nil {
		log.Println("Error while updating Repository:", err)
		return "", err
	}

	return repository.ID, nil
}

func (s *repositoryRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE repositories SET deleted_at = extract(epoch from current_timestamp) WHERE id = $1`

	_, err := s.DB.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error while deleting Repository :", err)
		return err
	}

	return nil
}

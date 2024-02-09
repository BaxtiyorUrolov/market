package postgres

import (
	"context"
	"fmt"
	"log"
	"market/api/models"
	"market/storage"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type staffTarifRepo struct {
	DB *pgxpool.Pool
}

func NewStaffTarifRepo(DB *pgxpool.Pool) storage.IStaffTariffRepo {
	return &staffTarifRepo{
		DB: DB,
	}
}

func (s *staffTarifRepo) Create(ctx context.Context, tarif models.CreateStaffTarif) (string, error) {
	id := uuid.New().String()
	createdAt := time.Now()

	if _, err := s.DB.Exec(ctx, `INSERT INTO staff_tarifs 
	(id, name, tarif_type, amount_for_cash, amount_for_card, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
			id, 
			tarif.Name, 
			tarif.TarifType, 
			tarif.AmountForCash, 
			tarif.AmountForCard, 
			createdAt,
	); err != nil {
		log.Println("Error while inserting data:", err)
		return "", err
	}

	return id, nil
}

func (s *staffTarifRepo) GetStaffTariffByID(ctx context.Context, id models.PrimaryKey) (models.StaffTarif, error) {
	staffTarif := models.StaffTarif{}
	query := `SELECT id, name, tarif_type, amount_for_cash, amount_for_card, created_at, updated_at, deleted_at FROM staff_tarifs WHERE id = $1`
	err := s.DB.QueryRow(ctx, query, id.ID).Scan(
		&staffTarif.ID,
		&staffTarif.Name,
		&staffTarif.TarifType,
		&staffTarif.AmountForCash,
		&staffTarif.AmountForCard,
		&staffTarif.CreatedAt,
		&staffTarif.UpdatedAt,
		&staffTarif.DeletedAt,
	)
	if err != nil {
		log.Println("Error while selecting staff tariff by ID:", err)
		return models.StaffTarif{}, err
	}
	return staffTarif, nil
}


func (s *staffTarifRepo) GetStaffTariffList(ctx context.Context, request models.GetListRequest) (models.StaffTarifResponse, error) {
	var (
		staffTarifs = []models.StaffTarif{}
		count       int
	)

	countQuery := `SELECT COUNT(*) FROM staff_tarifs`
	if request.Search != "" {
		countQuery += fmt.Sprintf(` WHERE name ILIKE '%%%s%%'`, request.Search)
	}

	err := s.DB.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		log.Println("Error while scanning count of staff tariffs:", err)
		return models.StaffTarifResponse{}, err
	}

	query := `SELECT id, name, tarif_type, amount_for_cash, amount_for_card, created_at, updated_at, deleted_at FROM staff_tarifs where deleted_at is null`
	if request.Search != "" {
		query += fmt.Sprintf(` WHERE name ILIKE '%%%s%%'`, request.Search)
	}
	query += ` LIMIT $1 OFFSET $2`

	rows, err := s.DB.Query(ctx, query, request.Limit, (request.Page-1)*request.Limit)
	if err != nil {
		log.Println("Error while querying staff tariffs:", err)
		return models.StaffTarifResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		staffTarif := models.StaffTarif{}
		err := rows.Scan(
			&staffTarif.ID,
			&staffTarif.Name,
			&staffTarif.TarifType,
			&staffTarif.AmountForCash,
			&staffTarif.AmountForCard,
			&staffTarif.CreatedAt,
			&staffTarif.UpdatedAt,
			&staffTarif.DeletedAt,
		)
		if err != nil {
			log.Println("Error while scanning row of staff tariffs:", err)
			return models.StaffTarifResponse{}, err
		}
		staffTarifs = append(staffTarifs, staffTarif)
	}

	return models.StaffTarifResponse{
		StaffTarifs: staffTarifs,
		Count:       count,
	}, nil
}

func (s *staffTarifRepo) UpdateStaffTariff(ctx context.Context, starif models.UpdateStaffTarif) (string, error) {
	query := `UPDATE staff_tarifs SET name = $1, tarif_type = $2, amount_for_cash = $3, amount_for_card = $4, updated_at = NOW() WHERE id = $5`

	_, err := s.DB.Exec(ctx, query,
		starif.Name,
		starif.TarifType,
		starif.AmountForCash,
		starif.AmountForCard,
		starif.ID,
	)
	if err != nil {
		log.Println("Error while updating Staff Tarif:", err)
		return "", err
	}

	return starif.ID, nil
}

func (s *staffTarifRepo) DeleteStaffTariff(ctx context.Context, id string) error {
	query := `UPDATE staff_tarifs SET deleted_at = NOW() WHERE id = $1`

	_, err := s.DB.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error while deleting Staff Tarif:", err)
		return err
	}

	return nil
}

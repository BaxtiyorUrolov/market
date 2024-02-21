package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"market/api/models"
	"market/pkg/logger"
	"market/storage"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type staffRepo struct {
	DB *pgxpool.Pool
	log logger.ILogger
}

func NewStaffRepo(DB *pgxpool.Pool, log logger.ILogger) storage.IStaffRepo {
	return &staffRepo{
		DB: DB,
		log: log,
	}
}

func (s *staffRepo) Create(ctx context.Context, staff models.CreateStaff) (string, error) {
	id := uuid.New().String()

	birthDate, err := time.Parse("2006-01-02", staff.BirthDate)
	if err != nil {
		log.Println("Error parsing birth date:", err)
		return "", err
	}
	age := uint(time.Since(birthDate).Hours() / 24 / 365)

	if _, err := s.DB.Exec(ctx, `INSERT INTO staffs 
		(id, branch_id, tariff_id, staff_type, name, balance, age, birth_date, login, password)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		id,
		staff.BranchID,
		staff.TariffID,
		staff.StaffType,
		staff.Name,
		staff.Balance,
		age,
		birthDate,
		staff.Login,
		staff.Password,
	); err != nil {
		log.Println("Error while inserting data ", err)
		return "", err
	}
	return id, nil
}

func (s *staffRepo) StaffByID(ctx context.Context, id models.PrimaryKey) (models.Staff, error) {
	var updatedAt sql.NullTime
	staff := models.Staff{}
	query := `SELECT id, branch_id, tariff_id, staff_type, name, balance, age, birth_date, login, created_at, updated_at FROM staffs WHERE id = $1`

	err := s.DB.QueryRow(ctx, query, id.ID).Scan(
		&staff.ID,
		&staff.BranchID,
		&staff.TariffID,
		&staff.StaffType,
		&staff.Name,
		&staff.Balance,
		&staff.Age,
		&staff.BirthDate,
		&staff.Login,
		&staff.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		log.Println("Error while selecting staff by ID:", err)
		return models.Staff{}, err
	}

	if updatedAt.Valid {
		staff.UpdatedAt = updatedAt.Time
	}

	return staff, nil
}

func (s *staffRepo) GetStaffTList(ctx context.Context, request models.GetListRequest) (models.StaffsResponse, error) {
	var (
		staffs = []models.Staff{}
		count     int
		updatedAt sql.NullTime
	)

	countQuery := `SELECT COUNT(*) FROM staffs where deleted_at = 0`
	if request.Search != "" {
		countQuery += fmt.Sprintf(` AND name ILIKE '%%%s%%' or login ilike '%%%s%%'`, request.Search, request.Search)
	}

	err := s.DB.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		log.Println("Error while scanning count of staffs:", err)
		return models.StaffsResponse{}, err
	}

	query := `SELECT id, branch_id, tariff_id, staff_type, name, balance, age, birth_date, login, created_at, updated_at FROM staffs where deleted_at = 0`
	if request.Search != "" {
		query += fmt.Sprintf(` AND name ILIKE '%%%s%%' or login ilike '%%%s%%'`, request.Search, request.Search)
	}
	query += ` LIMIT $1 OFFSET $2`

	rows, err := s.DB.Query(ctx, query, request.Limit, (request.Page-1)*request.Limit)
	if err != nil {
		log.Println("Error while querying staff :", err)
		return models.StaffsResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		staff := models.Staff{}
		err := rows.Scan(
			&staff.ID,
			&staff.BranchID,
			&staff.TariffID,
			&staff.StaffType,
			&staff.Name,
			&staff.Balance,
			&staff.Age,
			&staff.BirthDate,
			&staff.Login,
			&staff.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			log.Println("Error while scanning row of staffs:", err)
			return models.StaffsResponse{}, err
		}

		if updatedAt.Valid {
			staff.UpdatedAt = updatedAt.Time
		}

		staffs = append(staffs, staff)
	}

	return models.StaffsResponse{
		Staffs: staffs,
		Count:  count,
	}, nil
}

func (s *staffRepo) UpdateStaff(ctx context.Context, staff models.UpdateStaff) (string, error) {
	query := `UPDATE staffs SET branch_id = $1, tariff_id = $2, staff_type = $3, name = $4, balance = $5, login = $6, updated_at = NOW() WHERE id = $7`

	_, err := s.DB.Exec(ctx, query,
		&staff.BranchID,
		&staff.TariffID,
		&staff.StaffType,
		&staff.Name,
		&staff.Balance,
		&staff.Login,
		staff.ID,
	)
	if err != nil {
		log.Println("Error while updating Staff :", err)
		return "", err
	}

	return staff.ID, nil
}

func (s *staffRepo) DeleteStaff(ctx context.Context, id string) error {
	query := `UPDATE staffs SET deleted_at = extract(epoch from current_timestamp) WHERE id = $1`

	_, err := s.DB.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error while deleting Staff :", err)
		return err
	}

	return nil
}

func (s *staffRepo) GetPassword(ctx context.Context, id string) (string, error) {
	password := ""

	query := `
		select password from staffs 
		                where  id = $1`

	if err := s.DB.QueryRow(ctx, query, id).Scan(&password); err != nil {
		fmt.Println("Error while scanning password from users", err.Error())
		return "", err
	}

	return password, nil
}

func (s *staffRepo) UpdatePassword(ctx context.Context, request models.UpdateStaffPassword) error {
	query := `
		update staffs 
				set password = $1
					where id = $2`

	if _, err := s.DB.Exec(ctx, query, request.NewPassword, request.ID); err != nil {
		fmt.Println("error while updating password for staff", err.Error())
		return err
	}

	return nil
}

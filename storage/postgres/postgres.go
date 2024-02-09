package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"market/config"
	"market/storage"
)

type Store struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config) (storage.IStorage, error) {
	poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB))
	if err != nil {
		fmt.Println("error is while parsing config", err.Error())
		return nil, err
	}
	poolConfig.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		fmt.Println("error is while connecting to db", err.Error())
		return nil, err
	}
	return &Store{
		Pool: pool,
	}, nil
}

func (s *Store) Close() {
	s.Pool.Close()
}

func (s *Store) StaffTariff() storage.IStaffTariffRepo {
	return NewStaffTarifRepo(s.Pool)
}

func (s *Store) Category() storage.ICategory {
	return NewCategoryRepo(s.Pool)
}

func (s *Store) Product() storage.IProducts {
	return NewProductRepo(s.Pool)
}

func (s *Store) Branch() storage.IBranchStorage {
	return NewBranchRepo(s.Pool)
}

func (s *Store) Sale() storage.ISaleStorage {
	return NewSaleRepo(s.Pool)
}

func (s *Store) Transaction() storage.ITransactionStorage {
	return NewTransactionRepo(s.Pool)

}

func (s *Store) Staff() storage.IStaffRepo {
	return NewStaffRepo(s.Pool)
}

func (s *Store) Repository() storage.IRepositoryRepo {
	return NewRepositoryRepo(s.Pool)
}

func (s *Store) Basket() storage.IBasketRepo {
	return NewBasketRepo(s.Pool)
}

func (s *Store) RTransaction() storage.IRepositoryTransactionRepo {
	return NewRepositoryTransactionRepo(s.Pool)
}

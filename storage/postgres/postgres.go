package postgres

import (
	"context"
	"fmt"
	"market/config"
	"market/pkg/logger"
	"market/storage"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database"          //database is needed for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //postgres is used for database
	_ "github.com/golang-migrate/migrate/v4/source/file"       //file is needed for migration url
	_ "github.com/lib/pq"
)

type Store struct {
	Pool  *pgxpool.Pool
	log   logger.ILogger
	cfg   config.Config
}

func New(ctx context.Context, cfg config.Config, log logger.ILogger) (storage.IStorage, error) {
	url := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error("error while parsing config", logger.Error(err))
		return nil, err
	}

	poolConfig.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("error while connecting to db", logger.Error(err))
		return nil, err
	}

	//migration
	m, err := migrate.New("file://migrations/postgres/", url)
	if err != nil {
		log.Error("error while migrating", logger.Error(err))
		return nil, err
	}

	log.Info("???? came")

	if err = m.Up(); err != nil {
		log.Warning("migration up", logger.Error(err))
		if !strings.Contains(err.Error(), "no change") {
			fmt.Println("entered")
			version, dirty, err := m.Version()
			log.Info("version and dirty", logger.Any("version", version), logger.Any("dirty", dirty))
			if err != nil {
				log.Error("err in checking version and dirty", logger.Error(err))
				return nil, err
			}

			if dirty {
				version--
				if err = m.Force(int(version)); err != nil {
					log.Error("ERR in making force", logger.Error(err))
					return nil, err
				}
			}
			log.Warning("WARNING in migrating", logger.Error(err))
			return nil, err
		}
	}

	log.Info("!!!!! came here")

	return &Store{
		Pool:  pool,
		log:   log,
		cfg:   cfg,
	}, nil
}

func (s *Store) Close() {
	s.Pool.Close()
}

func (s *Store) StaffTariff() storage.IStaffTariffRepo {
	return NewStaffTarifRepo(s.Pool, s.log)
}

func (s *Store) Category() storage.ICategory {
	return NewCategoryRepo(s.Pool, s.log)
}

func (s *Store) Product() storage.IProducts {
	return NewProductRepo(s.Pool, s.log)
}

func (s *Store) Branch() storage.IBranchStorage {
	return NewBranchRepo(s.Pool, s.log)
}

func (s *Store) Sale() storage.ISaleStorage {
	return NewSaleRepo(s.Pool, s.log)
}

func (s *Store) Transaction() storage.ITransactionStorage {
	return NewTransactionRepo(s.Pool, s.log)

}

func (s *Store) Staff() storage.IStaffRepo {
	return NewStaffRepo(s.Pool, s.log)
}

func (s *Store) Repository() storage.IRepositoryRepo {
	return NewRepositoryRepo(s.Pool, s.log)
}

func (s *Store) Basket() storage.IBasketRepo {
	return NewBasketRepo(s.Pool, s.log)
}

func (s *Store) RTransaction() storage.IRepositoryTransactionRepo {
	return NewRepositoryTransactionRepo(s.Pool, s.log)
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	qrCreateCryptoExRate = `INSERT INTO crypto_exchange_rate(crypto, exchange_rate, date) VALUES ($1, $2, CURRENT_TIMESTAMP);`
	qrCryptoExRate       = `SELECT exchange_rate FROM crypto_exchange_rate WHERE crypto = $1 AND date < $2 ORDER BY date DESC LIMIT 1;`
)

type Storage struct {
	db *sql.DB
}

func New(storageDriver, storageInfo string) (*Storage, error) {
	const op = "internal.storage.postgres.New"

	db, err := sql.Open(storageDriver, storageInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) MigrationUp(storageURL, migrationPath string) error {
	const op = "storage.postgres.MigrationUp"

	migration, err := migrate.New(migrationPath, storageURL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = migration.Up()
	if err != nil && migration.Up().Error() != "no change" {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) MigrationDown(storageURL, migrationPath string) error {
	const op = "storage.postgres.MigrationDown"

	migration, err := migrate.New(migrationPath, storageURL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = migration.Down()
	if err != nil && migration.Down().Error() != "no change" {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddCryptoExRate(ctx context.Context, coin, price string) error {
	const op = "internal.storage.postgres.AddCryptoExRate"

	_, err := s.db.ExecContext(ctx, qrCreateCryptoExRate, coin, price)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetPrice(ctx context.Context, coin string, timestamp int) (string, error) {
	const op = "internal.storage.postgres.GetPrice"

	var price string

	timeFilter := time.Unix(int64(timestamp), 0)
	err := s.db.QueryRowContext(ctx, qrCryptoExRate, coin, timeFilter).Scan(&price)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return price, nil
}

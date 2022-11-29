package database

import (
	"errors"
	"github.com/uptrace/bun"
	"time"

	"github.com/huyhvq/ecommerce/assets"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun/dialect/pgdialect"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	*bun.DB
}

func New(driver string, dsn string, autoMigrate bool) (*DB, error) {
	sqlDB, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, err
	}
	db := bun.NewDB(sqlDB.DB, pgdialect.New())
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)
	if autoMigrate {
		iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
		if err != nil {
			return nil, err
		}
		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, dsn)
		if err != nil {
			return nil, err
		}
		defer migrator.Close()

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, err
		}
	}
	bdb := bun.NewDB(db.DB, pgdialect.New())
	return &DB{bdb}, nil
}

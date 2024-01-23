package db

import (
	"context"
	"database/sql"
	"emsrv/pkg/embedlog"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"hash/crc64"
)

// DB stores db connection
type DB struct {
	*pg.DB
	embedlog.Logger

	crcTable *crc64.Table
}

// New is a function that returns DB as wrapper on postgres connection.
func New(db *pg.DB) DB {
	d := DB{DB: db, crcTable: crc64.MakeTable(crc64.ECMA)}

	d.SetStdLoggers(true)
	return d
}

// Version is a function that returns Postgres version.
func (db *DB) Version() (string, error) {
	var v string
	if _, err := db.QueryOne(pg.Scan(&v), "select version()"); err != nil {
		return "", err
	}

	return v, nil
}

// runInTransaction runs chain of functions in transaction until first error
func (db *DB) runInTransaction(ctx context.Context, fns ...func(*pg.Tx) error) error {
	return db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		for _, fn := range fns {
			if err := fn(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *DB) RunMigration() error {
	dbConn, err := sql.Open("postgres", db.DB.Options().ToURL())
	if err != nil {
		return fmt.Errorf("[RunMigration] failed to open database connection: %w", err)
	}
	defer dbConn.Close()

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("[RunMigration] failed to create database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+"./docs",
		"", driver)
	if err != nil {
		return fmt.Errorf("[RunMigration] failed to create migration instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("[RunMigration] failed to apply migrations: %w", err)
	}

	return nil
}

/*// buildQuery applies all functions to orm query.
func buildQuery(ctx context.Context, db orm.DB, model interface{}, search Searcher, filters []Filter, pager Pager, ops ...OpFunc) *orm.Query {
	q := db.ModelContext(ctx, model)
	for _, filter := range filters {
		filter.Apply(q)
	}

	if reflect.ValueOf(search).IsValid() && !reflect.ValueOf(search).IsNil() { // is it good?
		search.Apply(q)
	}

	q = pager.Apply(q)
	applyOps(q, ops...)

	return q
}*/

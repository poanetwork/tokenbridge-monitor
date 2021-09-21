package db

import (
	"amb-monitor/config"
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	cfg     *config.DBConfig
	db      *sqlx.DB
}

func (db *DB) Migrate() error {
	m, err := migrate.New("file://db/migrations", db.dbURL("pgx"))
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (db *DB) dbURL(prefix string) string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s", prefix, db.cfg.User, db.cfg.Password, db.cfg.Host, db.cfg.Port, db.cfg.DB)
}

func NewDB(cfg *config.DBConfig) (*DB, error) {
	db := &DB{
		cfg:     cfg,
	}
	conn, err := sqlx.ConnectContext(context.Background(), "pgx", db.dbURL("postgres"))
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(3)
	conn.SetMaxOpenConns(10)
	db.db = conn
	return db, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	defer ObserveDuration(getCurrentFuncName(2))()
	return db.db.ExecContext(ctx, query, args...)
}

func (db *DB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	defer ObserveDuration(getCurrentFuncName(2))()
	return db.db.GetContext(ctx, dest, query, args...)
}

func (db *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	defer ObserveDuration(getCurrentFuncName(2))()
	return db.db.SelectContext(ctx, dest, query, args...)
}

func getCurrentFuncName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	details := runtime.FuncForPC(pc)
	if details == nil {
		return "unknown"
	}
	name := details.Name()
	name = name[strings.Index(name, ".")+1:]
	name = strings.TrimPrefix(name, "(*")
	name = strings.Replace(name, ")", "", 1)
	return name
}

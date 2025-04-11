package migrations

import (
	"embed"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed *.sql
var embedMigrations embed.FS

func SetupPostgres(pool *pgxpool.Pool, logger *zap.Logger) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("cannot set dialect in goose", zap.Error(err))
		os.Exit(-1)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	if err := goose.Up(db, "."); err != nil {
		logger.Error("cannot apply migrations", zap.Error(err))
		os.Exit(-1)
	}

	logger.Info("Postgres migrations applied successfully")
}

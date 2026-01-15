package db

import (
	"context"
	"fmt"
	"gin-user-management/internal/config"
	"gin-user-management/internal/db/sqlc"
	"gin-user-management/internal/util"
	"gin-user-management/pkg/pgx"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

var (
	DB     sqlc.Querier
	DBPool *pgxpool.Pool
)

func InitDB() error {
	connStr := config.NewConfig().DNS()
	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("Error parsing DB config: %v", err)
	}

	sqlLogger := util.NewLogger("logs/sql.log", "info")

	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger: &pgx.PgxZerologTracer{
			Logger:         *sqlLogger,
			SlowQueryLimit: 500 * time.Millisecond,
		},
		LogLevel: tracelog.LogLevelDebug,
	}
	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DBPool, err = pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("Error create DB pool: %v", err)
	}
	DB = sqlc.New(DBPool)

	if err := DBPool.Ping(ctx); err != nil {
		return fmt.Errorf("DB ping error: %v", err)
	}

	log.Println("🐘 Database connected.")

	return nil
}

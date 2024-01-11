package repository

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
	cache "kiramishima/m-backend/internal/adapters/cache/redis"
	"kiramishima/m-backend/internal/core/domain"
	"time"
)

var DatabaseModule = fx.Module("db",
	fx.Provide(NewDatabase),
	fx.Provide(func(conn *sqlx.DB) *AuthRepository {
		return NewAuthRepository(conn)
	}),
	fx.Provide(func(conn *sqlx.DB, cache *cache.RedisCache) *BondRepository {
		return NewBondRepository(conn, cache)
	}),
	fx.Provide(func(conn *sqlx.DB, cache *cache.RedisCache) *MarketBondRepository {
		return NewMarketBondRepository(conn, cache)
	}),
	fx.Provide(func(conn *sqlx.DB, cache *cache.RedisCache) *UserRepository {
		return NewUserRepository(conn, cache)
	}),
)

// NewDatabase creates an instance of DB
func NewDatabase(lc fx.Lifecycle, cfg *domain.Configuration, logger *zap.SugaredLogger) (*sqlx.DB, error) {

	db, err := sqlx.Connect(cfg.DatabaseDriver, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	// conf connections
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	// iddletime
	duration, err := time.ParseDuration(cfg.MaxIdleTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	// context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ContextTimeout)*time.Second)
	defer cancel()

	// Ping to DB
	status := "up"
	err = db.PingContext(ctx)
	if err != nil {
		status = "down"
		return nil, err
	}
	logger.Debugf("Status DB: %s", status)
	return db, nil
}

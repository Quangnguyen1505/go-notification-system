package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	configs "github.com/quangnguyen1505/go-notification-system/pkg/config"
	"github.com/quangnguyen1505/go-notification-system/pkg/logger"
	"go.uber.org/zap"
)

const (
	_defaultConnAttempts = 3
	_defaultConnTimeout  = time.Second * 5
)

type DBConnString string

type postgres struct {
	connAttempts int
	connTimeout  time.Duration

	db *pgxpool.Pool
}

var _ DBEngine = (*postgres)(nil)

func NewPostgresDB(config *configs.Postgres, logger *logger.LoggerZap) (DBEngine, error) {
	baseCtx := context.Background()
	logger.Info("m.Host, m.Username, m.Password, m.Dbname, m.Port::", zap.String("host", config.Host), zap.String("username", config.Username), zap.String("password", config.Password), zap.String("dbname", config.Dbname), zap.Int("port", config.Port))
	dsn := "user=%s password=%s dbname=%s host=%s port=%d sslmode=disable"
	var s = fmt.Sprintf(dsn, config.Username, config.Password, config.Dbname, config.Host, config.Port)

	ctx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
	defer cancel()

	cfgPool, err := pgxpool.ParseConfig(s)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	maxConns := config.MaxOpenConns
	if maxConns < 1 {
		maxConns = 1
	}
	minConns := config.MinConns
	if minConns < 0 {
		minConns = 0
	}
	if minConns > maxConns {
		minConns = maxConns
	}

	cfgPool.MaxConns = int32(maxConns)
	cfgPool.MinConns = int32(minConns)
	cfgPool.MaxConnLifetime = time.Duration(config.ConnMaxLifeTime) * time.Second

	pg := &postgres{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	pg.db, err = pgxpool.NewWithConfig(ctx, cfgPool)
	if err != nil {
		return nil, err
	}

	for i := 0; i < pg.connAttempts; i++ {
		pingCtx, pingCancel := context.WithTimeout(baseCtx, 5*time.Second)
		err = pg.db.Ping(pingCtx)
		pingCancel()
		if err == nil {
			break
		}

		logger.Warn("retry connect", zap.Int("attempt", i+1), zap.Error(err))
		if i < pg.connAttempts-1 {
			time.Sleep(pg.connTimeout)
		}
	}

	if err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to PostgreSQL")

	return pg, nil
}

func (p *postgres) Configure(opts ...Option) DBEngine {
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *postgres) GetDB() *pgxpool.Pool {
	return p.db
}

func (p *postgres) Close() {
	if p.db != nil {
		p.db.Close()
	}
}

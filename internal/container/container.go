package container

import (
	"context"
	"log/slog"
	"os"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/Beriw98/user-management/ent"
	"github.com/Beriw98/user-management/internal/config"
	"github.com/Beriw98/user-management/internal/infrastructure/database/repository"
	"github.com/Beriw98/user-management/internal/infrastructure/httpsrv/handler"
)

type Container struct {
	DB             *ent.Client
	Logger         *slog.Logger
	UserRepository *repository.User
	UserHandler    *handler.UserHTTPHandler
}

func NewContainer(cfg *config.Config) (*Container, error) {
	ctx := context.Background()
	pgPoolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}

	pgPoolCfg.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, pgPoolCfg)
	if err != nil {
		return nil, err
	}
	drv := sql.OpenDB(dialect.Postgres, stdlib.OpenDB(*pool.Config().ConnConfig))
	client := ent.NewClient(ent.Driver(drv))

	userRepository := repository.NewUserRepository(client)
	userHandler := handler.NewUserHTTPHandler(userRepository)

	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(l)

	return &Container{
		DB:             client,
		UserRepository: userRepository,
		UserHandler:    userHandler,
		Logger:         l,
	}, nil
}

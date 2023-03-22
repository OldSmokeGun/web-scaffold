//go:build wireinject
// +build wireinject

package command

import (
	"context"
	"database/sql"

	"go-scaffold/internal/app"
	"go-scaffold/internal/app/adapter/cron"
	"go-scaffold/internal/app/adapter/server"
	"go-scaffold/internal/app/pkg"
	"go-scaffold/internal/config"
	"go-scaffold/pkg/trace"

	"github.com/google/wire"
	"golang.org/x/exp/slog"
)

func initServer(
	context.Context,
	config.AppName,
	config.Env,
	*slog.Logger,
	*trace.Trace,
) (*server.Server, func(), error) {
	panic(wire.Build(
		config.ProviderSet,
		app.ProviderSet,
	))
}

func initCron(
	context.Context,
	config.AppName,
	config.Env,
	*slog.Logger,
) (*cron.Cron, func(), error) {
	panic(wire.Build(
		config.ProviderSet,
		app.ProviderSet,
	))
}

func initDB(
	context.Context,
	config.DBConn,
	*slog.Logger,
) (*sql.DB, func(), error) {
	panic(wire.Build(
		pkg.ProviderSet,
	))
}

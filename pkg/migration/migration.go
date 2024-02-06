package migration

import (
	"os"

	"github.com/Dostonlv/todo.git/pkg/config"
	"github.com/Dostonlv/todo.git/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Invoke(New),
)

type Params struct {
	fx.In
	Logger logger.ILogger
	Config config.IConfig
}

func New(p Params) {
	m, err := migrate.New("file://migrations", p.Config.GetString("database.migration"))
	if err != nil {
		p.Logger.Error("err from migration.New", zap.Error(err))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		p.Logger.Error("err from up migration", zap.Error(err))
		os.Exit(1)
	}
}

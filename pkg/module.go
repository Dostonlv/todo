package pkg

import (
	"github.com/Dostonlv/todo.git/pkg/casbin"
	"github.com/Dostonlv/todo.git/pkg/db"
	"github.com/Dostonlv/todo.git/pkg/logger"
	"github.com/Dostonlv/todo.git/pkg/migration"
	"go.uber.org/fx"
)

var Module = fx.Options(
	casbin.Module,
	migration.Module,
	db.Module,
	logger.Module,
)

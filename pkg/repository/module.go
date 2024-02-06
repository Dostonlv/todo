package repository

import (
	"github.com/Dostonlv/todo.git/pkg/repository/postgres"
	"go.uber.org/fx"
)

var Module = fx.Options(
	postgres.Module,
)

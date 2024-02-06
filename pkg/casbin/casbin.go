package casbin

import (
	"github.com/Dostonlv/todo.git/pkg/config"
	"github.com/Dostonlv/todo.git/pkg/logger"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(New)

type AuthEnforcer interface {
	LoadPolicy() *casbin.Enforcer
}

type Params struct {
	fx.In
	Logger logger.ILogger
	Config config.IConfig
}

type module struct {
	logger       logger.ILogger
	config       config.IConfig
	authEnforcer *casbin.Enforcer
}

func New(p Params) (AuthEnforcer, error) {
	var (
		dns = p.Config.GetString("database.casbin")
		err error
	)

	adapter, err := gormadapter.NewAdapter("postgres", dns, true)
	if err != nil {
		p.Logger.Error("err on gormadapter.NewAdapter: %v", zap.Error(err))
		return nil, err
	}

	authEnforcer, err := casbin.NewEnforcer("./configs/auth_model.conf", adapter)
	if err != nil {
		p.Logger.Error("err on casbin.NewEnforcer: %v", zap.Error(err))
		return nil, err
	}

	err = authEnforcer.LoadPolicy()
	if err != nil {
		p.Logger.Error("err on authEnforcer.LoadPolicy: %v", zap.Error(err))
		return nil, err
	}

	return &module{
		logger:       p.Logger,
		config:       p.Config,
		authEnforcer: authEnforcer,
	}, nil
}

func (m *module) LoadPolicy() *casbin.Enforcer {
	return m.authEnforcer
}

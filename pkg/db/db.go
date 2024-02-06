package db

import (
	"context"

	"go.uber.org/fx"

	"github.com/Dostonlv/todo.git/pkg/config"
	"github.com/Dostonlv/todo.git/pkg/logger"
	"github.com/Dostonlv/todo.git/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Module = fx.Options(
	fx.Provide(NewDBConn),
)

type Querier interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}

type Params struct {
	fx.In
	Config config.IConfig
	Logger logger.ILogger
}

type dbConn struct {
	config config.IConfig
	dbPool *pgxpool.Pool
	logger logger.ILogger
}

func NewDBConn(params Params) (Querier, error) {

	var (
		dns = params.Config.GetString("database.dns")
		err error
	)

	db, err := pgxpool.New(context.Background(), dns)
	if err != nil {
		params.Logger.Error("Err on pgxpool.Connect(%v): %v", dns, err)
		return nil, err
	}

	return &dbConn{
		dbPool: db,
		logger: params.Logger,
		config: params.Config,
	}, nil
}

func (db *dbConn) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	db.logger.Info("DB: Exec sql: [%s] Args: [%v]", utils.RemoveSpecialChars2(sql), args)
	return db.dbPool.Exec(ctx, sql, args...)
}

func (db *dbConn) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	db.logger.Info("DB: Query sql: [%s] Args: [%v]", utils.RemoveSpecialChars2(sql), args)
	return db.dbPool.Query(ctx, sql, args...)
}

func (db *dbConn) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	db.logger.Info("DB: QueryRow sql: [%s] Args: [%v]", utils.RemoveSpecialChars2(sql), args)
	return db.dbPool.QueryRow(ctx, sql, args...)
}

func (db *dbConn) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.dbPool.Begin(ctx)
}

func (db *dbConn) SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults {
	return db.dbPool.SendBatch(ctx, batch)
}

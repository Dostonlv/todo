package userRepo

import (
	"context"
	"errors"
	"time"

	"github.com/Dostonlv/todo.git/internal/structs"
	"github.com/Dostonlv/todo.git/pkg/db"
	"github.com/Dostonlv/todo.git/pkg/logger"
	"github.com/Dostonlv/todo.git/pkg/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	Module  = fx.Provide(New)
	logPath = "pkg/repository/postgres/users_repo/users.go"
)

type (
	Params struct {
		fx.In
		Logger logger.ILogger
		DB     db.Querier
	}

	Repo interface {
		Create(ctx context.Context, user structs.User) error
		GetUserByID(ctx context.Context, id int) (user structs.User, err error)
		GetUsers(ctx context.Context, request structs.Filter) (structs.UserList, error)
		Update(ctx context.Context, user structs.User) error
		Delete(ctx context.Context, id int) error
	}
	repo struct {
		logger logger.ILogger
		db     db.Querier
	}
)

func New(p Params) Repo {
	return &repo{
		logger: p.Logger,
		db:     p.DB,
	}
}

func (r *repo) Create(ctx context.Context, user structs.User) error {
	logPath := logPath + "->Create"

	var pgErr = &pgconn.PgError{}
	exec, err := r.db.Exec(ctx, `
	insert into users(
		login,
		password
	) values ($1, $2)`,
		user.Login,
		user.Password)
	if err != nil {
		errors.As(err, &pgErr)
		if pgerrcode.UniqueViolation == pgErr.Code {
			return structs.ErrUniqueViolation
		}
		r.logger.Error(logPath+" err on r.db.Exec", zap.Error(err))
		return err
	}
	if exec.RowsAffected() == 0 {
		return structs.ErrNoRowsAffected
	}

	return err

}

func (r repo) GetUserByID(ctx context.Context, id int) (user structs.User, err error) {

	r.logger.Info("GetUserByID", zap.Int("id", id))

	return r.selectUser(ctx, `u.id = $1`, id)
}

func (r repo) selectUser(ctx context.Context, c string, v ...interface{}) (user structs.User, err error) {
	logPath := logPath + "->selectUser"

	var (
		createdAt time.Time
	)

	err = r.db.QueryRow(ctx, `
	select
		u.id,
		u.login,
		u.password,
		u.created_at
	from users u
	where `+c, v...).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, structs.ErrNotFound
		}
		r.logger.Error(logPath+" err from r.db.QueryRow", zap.Error(err))
		return user, err
	}

	user.CreatedAt = createdAt.Format(utils.Layout)

	return user, err

}

func (r repo) GetUsers(ctx context.Context, request structs.Filter) (structs.UserList, error) {
	logPath := logPath + "->GetUsers"

	var (
		users     []structs.User
		createdAt time.Time
	)

	rows, err := r.db.Query(ctx, `
	select
		u.id,
		u.login,
		u.password,
		u.created_at
	from users u
	where u.login ilike $1
	order by u.id
	offset $2
	limit $3`, request.Search, request.Offset, request.Limit)
	if err != nil {
		r.logger.Error(logPath+" err from r.db.Query", zap.Error(err))
		return structs.UserList{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.User
		err = rows.Scan(
			&user.ID,
			&user.Login,
			&user.Password,
			&createdAt,
		)
		if err != nil {
			r.logger.Error(logPath+" err from rows.Scan", zap.Error(err))
			return structs.UserList{}, err
		}
		user.CreatedAt = createdAt.Format(utils.Layout)
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(logPath+" err from rows.Err", zap.Error(err))
		return structs.UserList{}, err
	}

	return structs.UserList{Users: users, Count: len(users)}, nil
}

func (r repo) Update(ctx context.Context, user structs.User) error {
	logPath := logPath + "->Update"

	var pgErr = &pgconn.PgError{}
	exec, err := r.db.Exec(ctx, `
	update users
	set
		login = $1,
		password = $2
	where id = $3`,
		user.Login,
		user.Password,
		user.ID)
	if err != nil {
		errors.As(err, &pgErr)
		if pgerrcode.UniqueViolation == pgErr.Code {
			return structs.ErrUniqueViolation
		}
		r.logger.Error(logPath+" err on r.db.Exec", zap.Error(err))
		return err
	}
	if exec.RowsAffected() == 0 {
		return structs.ErrNoRowsAffected
	}

	return err
}

func (r repo) Delete(ctx context.Context, id int) error {
	logPath := logPath + "->Delete"

	exec, err := r.db.Exec(ctx, `
	delete from users
	where id = $1`, id)
	if err != nil {
		r.logger.Error(logPath+" err on r.db.Exec", zap.Error(err))
		return err
	}
	if exec.RowsAffected() == 0 {
		return structs.ErrNoRowsAffected
	}

	return err
}

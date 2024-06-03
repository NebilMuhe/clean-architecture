package user

import (
	db "clean-architecture/internal/constants/db/sqlc"
	"clean-architecture/internal/constants/dbinstance"
	"clean-architecture/internal/constants/errors"
	"clean-architecture/internal/constants/model/usermodel"
	"clean-architecture/internal/data"
	"clean-architecture/utils/logger"
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type user struct {
	db  dbinstance.DBInstance
	log logger.Logger
}

func Init(db dbinstance.DBInstance, log logger.Logger) data.User {
	return &user{
		db:  db,
		log: log,
	}
}

func (u *user) CreateUser(ctx context.Context, param usermodel.RegisterUser) (*usermodel.User, error) {
	arg := db.CreateUserParams{
		Username: param.Username,
		Email:    param.Email,
		Password: param.Password,
	}
	user, err := u.db.CreateUser(ctx, arg)
	if err != nil {
		u.log.Error(ctx, "unable to create user", zap.Error(err))
		err := errors.ErrWriteError.Wrap(err, "unable to create user")
		return nil, err
	}

	return &usermodel.User{
		ID:        user.ID.Bytes,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (u *user) LoginUser(ctx context.Context, param usermodel.LoginUser) (*usermodel.User, error) {
	usr, err := u.db.LoginUser(ctx, param.Username)
	if err != nil {
		u.log.Error(ctx, "unable to login", zap.Error(err), zap.String("username", param.Username))
		if err == pgx.ErrNoRows {
			err = errors.ErrNoRecordFound.Wrap(err, "user does not exists")
		} else {
			errors.ErrReadError.Wrap(err, "unable to login")
		}

		return nil, err
	}

	return &usermodel.User{
		ID:       usr.ID.Bytes,
		Username: usr.Username,
		Password: usr.Password,
	}, nil
}

func (u *user) RefreshToken(ctx context.Context, username string, refToken string) (*usermodel.RefreshToken, error) {
	arg := db.CreateSessionParams{
		Username:     username,
		RefreshToken: refToken,
	}
	session, err := u.db.CreateSession(ctx, arg)
	if err != nil {
		u.log.Error(ctx, "unable to create session", zap.Error(err))
		err := errors.ErrWriteError.Wrap(err, "unable to create")
		return nil, err
	}
	return &usermodel.RefreshToken{
		ID:           session.ID.Bytes,
		Username:     session.Username,
		RefreshToken: session.RefreshToken,
		CreatedAt:    session.CreatedAt.Time,
	}, nil
}

func (u *user) GetRefreshToken(ctx context.Context, username string) (string, error) {
	session, err := u.db.GetToken(ctx, username)
	if err != nil {
		u.log.Error(ctx, "unable to read", zap.Error(err))
		err := errors.ErrReadError.Wrap(err, "unable to read")
		return "", err
	}

	return session.RefreshToken, nil
}

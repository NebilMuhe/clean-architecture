package user

import (
	db "clean-architecture/internal/constants/db/sqlc"
	"clean-architecture/internal/constants/dbinstance"
	"clean-architecture/internal/constants/errors"
	"clean-architecture/internal/constants/model/usermodel"
	"clean-architecture/internal/data"
	"clean-architecture/utils/logger"
	"context"

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

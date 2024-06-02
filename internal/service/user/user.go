package user

import (
	"clean-architecture/internal/constants/errors"
	"clean-architecture/internal/constants/model/usermodel"
	"clean-architecture/internal/data"
	"clean-architecture/internal/service"
	"clean-architecture/utils/helpers"
	"clean-architecture/utils/logger"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type user struct {
	data data.User
	log  logger.Logger
}

func Init(u data.User, log logger.Logger) service.User {
	return &user{
		data: u,
		log:  log,
	}
}

func (u *user) CreateUser(ctx context.Context, param usermodel.RegisterUser) (*usermodel.User, error) {
	if err := param.Validate(ctx, u.log); err != nil {
		return nil, err
	}

	password, err := helpers.HashPassword(param.Password)
	if err != nil {
		u.log.Error(ctx, "unable to hash password", zap.Error(err))
		err := errors.ErrWriteError.Wrap(err, "unable to hash password")
		return nil, err
	}

	param.Password = password

	user, err := u.data.CreateUser(ctx, param)
	if err != nil {
		return nil, err
	}

	fmt.Println("user", user)

	return user, nil
}

func (u *user) LoginUser(ctx context.Context, param usermodel.LoginUser) (map[string]string, error) {
	usr, err := u.data.LoginUser(ctx, param)
	if err != nil {
		return nil, err
	}

	err = helpers.CheckPassword(usr.Password, param.Password)
	if err != nil {
		u.log.Error(ctx, "invalid password", zap.Error(err))
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid password")
		return nil, err
	}

	token, err := helpers.CreateToken(usr.ID.String(), usr.Username)
	if err != nil {
		u.log.Error(ctx, "unable to create token", zap.Error(err))
		err = errors.ErrWriteError.Wrap(err, "unable to create token")
		return nil, err
	}

	return token, nil
}

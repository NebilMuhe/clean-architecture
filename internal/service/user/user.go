package user

import (
	"clean-architecture/internal/constants/errors"
	"clean-architecture/internal/constants/model/usermodel"
	"clean-architecture/internal/data"
	"clean-architecture/internal/service"
	"clean-architecture/utils/logger"
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

	password, err := HashPassword(param.Password)
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

	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
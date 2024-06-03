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

	"github.com/spf13/viper"
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

	password, err := helpers.HashPassword(ctx, param.Password, u.log)
	if err != nil {
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

func (u *user) LoginUser(ctx context.Context, param usermodel.LoginUser) (*usermodel.Token, error) {
	usr, err := u.data.LoginUser(ctx, param)
	if err != nil {
		return nil, err
	}

	err = helpers.CheckPassword(ctx, usr.Password, param.Password, u.log)
	if err != nil {
		return nil, err
	}

	token, err := helpers.CreateToken(ctx, usr.ID.String(), usr.Username, u.log)
	if err != nil {
		return nil, err
	}

	rfEncrypted, err := helpers.Encrypt(ctx, []byte(viper.GetString("secret_key")), token.RefreshToken, u.log)
	if err != nil {
		return nil, err
	}

	_, err = u.data.RefreshToken(ctx, param.Username, rfEncrypted)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (u *user) RefreshToken(ctx context.Context, tokenString string) (*usermodel.Token, error) {
	err := helpers.VerifyToken(ctx, tokenString, u.log)
	if err != nil {
		return nil, err
	}

	res, err := helpers.ExtractUsernameAndID(ctx, tokenString, u.log)
	if err != nil {
		u.log.Error(ctx, "unable to extract username and id", zap.Error(err))
		err := errors.ErrReadError.Wrap(err, "unable to read")
		return nil, err
	}

	rfToken, err := u.data.GetRefreshToken(ctx, res["username"])
	if err != nil {
		return nil, err
	}

	decrytpRfToken, err := helpers.Decrypt(ctx, []byte(viper.GetString("secret_key")), rfToken, u.log)
	if err != nil {
		return nil, err
	}

	if decrytpRfToken != tokenString {
		err := errors.ErrInvalidUserInput.New("invalid token")
		u.log.Error(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	token, err := helpers.CreateToken(ctx, res["id"], res["username"], u.log)
	if err != nil {
		u.log.Error(ctx, "unable to create token", zap.Error(err))
		err = errors.ErrWriteError.Wrap(err, "unable to create token")
		return nil, err
	}

	rfEncrypted, err := helpers.Encrypt(ctx, []byte(viper.GetString("secret_key")), token.RefreshToken, u.log)
	if err != nil {
		return nil, err
	}

	_, err = u.data.RefreshToken(ctx, res["username"], rfEncrypted)
	if err != nil {
		return nil, err
	}
	return token, nil
}

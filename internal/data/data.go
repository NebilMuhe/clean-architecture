package data

import (
	"clean-architecture/internal/constants/model/usermodel"
	"context"
)

type User interface {
	CreateUser(ctx context.Context, param usermodel.RegisterUser) (*usermodel.User, error)
	LoginUser(ctx context.Context, param usermodel.LoginUser) (*usermodel.User, error)
	RefreshToken(ctx context.Context, username string, refToken string) (*usermodel.RefreshToken, error)
	GetRefreshToken(ctx context.Context, username string) (string, error)
}

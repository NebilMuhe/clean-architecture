package service

import (
	"clean-architecture/internal/constants/model/usermodel"
	"context"
)

type User interface {
	CreateUser(ctx context.Context, param usermodel.RegisterUser) (*usermodel.User, error)
	LoginUser(ctx context.Context, param usermodel.LoginUser) (*usermodel.Token, error)
	RefreshToken(ctx context.Context, tokenString string) (*usermodel.Token, error)
}

package service

import (
	"clean-architecture/internal/constants/model/usermodel"
	"context"
)

type User interface {
	CreateUser(ctx context.Context, param usermodel.RegisterUser) (*usermodel.User, error)
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	LoginUser(ctx context.Context, username string) (User, error)
}

var _ Querier = (*Queries)(nil)

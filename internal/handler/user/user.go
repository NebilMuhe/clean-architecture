package user

import (
	"clean-architecture/internal/constants"
	"clean-architecture/internal/constants/errors"
	"clean-architecture/internal/constants/model/usermodel"
	"clean-architecture/internal/handler"
	"clean-architecture/internal/service"
	"clean-architecture/utils/logger"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type user struct {
	service service.User
	log     logger.Logger
	timeout time.Duration
}

func Init(service service.User, log logger.Logger, timeout time.Duration) handler.User {
	return &user{
		service: service,
		log:     log,
	}
}

func (u *user) CreateUser(ctx *gin.Context) {
	// contx, cancel := context.WithTimeout(ctx, u.timeout)
	// defer cancel()

	var usr usermodel.RegisterUser

	if err := ctx.ShouldBind(&usr); err != nil {
		u.log.Error(ctx, "unable to bind user data", zap.Error(err))
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid user input")
		ctx.Error(err)
		return
	}

	user, err := u.service.CreateUser(ctx, usr)
	if err != nil {
		ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusCreated, user, nil)
}

func (u *user) LoginUser(ctx *gin.Context) {
	var usr usermodel.LoginUser

	if err := ctx.ShouldBind(&usr); err != nil {
		u.log.Error(ctx, "unable to bind user data", zap.Error(err))
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid user input")
		ctx.Error(err)
		return
	}

	res, err := u.service.LoginUser(ctx, usr)
	if err != nil {
		ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, res, nil)
}

func (u *user) RefreshToken(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
		err := errors.ErrInvalidUserInput.New("invalid input")
		u.log.Error(ctx, "invalid token", zap.Error(err))
		ctx.Error(err)
		return
	}

	tokenString := authorization[len("Bearer "):]
	if tokenString == "" {
		err := errors.ErrInvalidUserInput.New("invalid input")
		u.log.Error(ctx, "invalid token", zap.Error(err))
		return
	}

	token, err := u.service.RefreshToken(ctx, tokenString)
	if err != nil {
		ctx.Error(err)
		return
	}

	constants.SuccessResponse(ctx, http.StatusOK, token, nil)
}

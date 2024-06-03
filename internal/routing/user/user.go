package user

import (
	"clean-architecture/internal/handler"
	"clean-architecture/internal/routing"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoute(grp *gin.RouterGroup, user handler.User) {
	users := grp.Group("users")
	userRoutes := []routing.Router{
		{
			Path:    "/register",
			Method:  http.MethodPost,
			Handler: user.CreateUser,
		},
		{
			Path:    "/login",
			Method:  http.MethodPost,
			Handler: user.LoginUser,
		},
		{
			Path:    "/refresh",
			Method:  http.MethodPost,
			Handler: user.RefreshToken,
		},
	}

	routing.RegisterRoutes(users, userRoutes)
}

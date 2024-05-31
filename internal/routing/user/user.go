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
			Method:  http.MethodPost,
			Handler: user.CreateUser,
		},
	}

	routing.RegisterRoutes(users, userRoutes)
}

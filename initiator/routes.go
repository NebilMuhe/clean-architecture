package initiator

import (
	"clean-architecture/internal/routing/user"

	"github.com/gin-gonic/gin"
)

func InitRoute(group *gin.RouterGroup, handler Handler) {
	user.InitRoute(group, handler.handler)
}

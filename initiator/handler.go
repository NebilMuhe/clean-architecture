package initiator

import (
	"clean-architecture/internal/handler"
	"clean-architecture/internal/handler/user"
	"clean-architecture/utils/logger"
	"time"
)

type Handler struct {
	handler handler.User
}

func InitHandler(service Service, log logger.Logger) Handler {
	return Handler{
		handler: user.Init(service.user, log.Named("handler-layer"), time.Minute*1),
	}
}

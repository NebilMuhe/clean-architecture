package initiator

import (
	"clean-architecture/internal/service"
	"clean-architecture/internal/service/user"
	"clean-architecture/utils/logger"
)

type Service struct {
	user service.User
}

func InitService(persitence Persistence, log logger.Logger) Service {
	return Service{
		user: user.Init(persitence.user, log.Named("service-layer")),
	}
}

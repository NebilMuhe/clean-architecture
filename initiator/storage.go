package initiator

import (
	"clean-architecture/internal/constants/dbinstance"
	"clean-architecture/internal/data"
	"clean-architecture/internal/data/user"
	"clean-architecture/utils/logger"
)

type Persistence struct {
	user data.User
}

func InitPersistence(db dbinstance.DBInstance, log logger.Logger) Persistence {
	return Persistence{
		user: user.Init(db, log.Named("persistence or data layer")),
	}
}

package initiator

import (
	"clean-architecture/utils/logger"
	"log"

	"go.uber.org/zap"
)

func InitLogger() logger.Logger {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	return logger.NewLogger(l)
}

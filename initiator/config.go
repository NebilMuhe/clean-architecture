package initiator

import (
	"clean-architecture/utils/logger"
	"context"

	"github.com/spf13/viper"
)

func InitConfig(ctx context.Context, name, path string, log logger.Logger) {
	viper.SetConfigName(name)
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(ctx, "unable to read configuration file")
	}
}

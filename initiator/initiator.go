package initiator

import (
	"context"

	"github.com/spf13/viper"
)

func Initialize(ctx context.Context) {
	log := InitLogger()
	log.Info(ctx, "logger initialized")

	log.Info(ctx, "initializing configuration")
	InitConfig(ctx, "config", "config", log)
	log.Info(ctx, "initialized configuration")

	log.Info(ctx, "initializing database")
	InitDB(ctx, viper.GetString("database.url"), log)
	log.Info(ctx, "initilaizied database")

	log.Info(ctx, "initializing migration")
	InitMigration(ctx, viper.GetString("database.file"), viper.GetString("database.murl"), log)
	log.Info(ctx, "initialized migration")
}

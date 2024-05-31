package initiator

import (
	"clean-architecture/internal/constants/dbinstance"
	"clean-architecture/internal/handler/middleware"
	"context"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Initialize(ctx context.Context) {
	log := InitLogger()
	log.Info(ctx, "logger initialized")

	log.Info(ctx, "initializing configuration")
	InitConfig(ctx, "config", "config", log)
	log.Info(ctx, "initialized configuration")

	log.Info(ctx, "initializing database")
	pool := InitDB(ctx, viper.GetString("database.url"), log)
	log.Info(ctx, "initilaizied database")

	log.Info(ctx, "initializing migration")
	InitMigration(ctx, viper.GetString("database.file"), viper.GetString("database.murl"), log)
	log.Info(ctx, "initialized migration")

	log.Info(ctx, "initializing persistence layer")
	persitence := InitPersistence(dbinstance.New(pool), log)
	log.Info(ctx, "initialized persistence layer")

	log.Info(ctx, "initializing service layer")
	service := InitService(persitence, log)
	log.Info(ctx, "initialized service layer")

	log.Info(ctx, "initializing handler layer")
	handler := InitHandler(service, log)
	log.Info(ctx, "initialized handler")

	log.Info(ctx, "intializing server")
	server := gin.New()
	server.Use(ginzap.RecoveryWithZap(log.GetZapLogger().Named("gin-recovery"), true))
	server.Use(middleware.ErrorHandler())
	log.Info(ctx, "initialized server")

	log.Info(ctx, "initializing routes")
	router := server.Group("/v1")
	InitRoute(router, handler)
	log.Info(ctx, "initialized routes")
}

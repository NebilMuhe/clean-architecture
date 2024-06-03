package initiator

import (
	"clean-architecture/internal/constants/dbinstance"
	"clean-architecture/internal/handler/middleware"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	router := server.Group("/api/v1")
	InitRoute(router, handler)
	log.Info(ctx, "initialized routes")

	srv := http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: server,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		log.Info(ctx, "started listening on server with ", zap.Int("port", viper.GetInt("server.port")))
		log.Info(ctx, fmt.Sprintf("server stopped with the error %v", srv.ListenAndServe()))
	}()

	sig := <-quit
	log.Info(ctx, fmt.Sprintf("server shutting down with signal %v", sig))
	cntx, cancel := context.WithTimeout(ctx, viper.GetDuration("server.timeout"))
	defer cancel()

	log.Info(ctx, "shutting down server")
	err := srv.Shutdown(cntx)

	if err != nil {
		log.Fatal(ctx, "error while shutting down server", zap.Error(err))
	}

	log.Info(ctx, "server shutdown complete")
}

package initiator

import "context"

func Initialize(ctx context.Context) {
	log := InitLogger()
	log.Info(ctx, "logger initialized")

	log.Info(ctx, "initializing configuration")
	InitConfig(ctx, "config", "config", log)
	log.Info(ctx, "initialized configuration")
}

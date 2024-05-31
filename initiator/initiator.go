package initiator

import "context"

func Initialize(ctx context.Context) {
	log := InitLogger()
	log.Info(ctx, "logger initialized")
}

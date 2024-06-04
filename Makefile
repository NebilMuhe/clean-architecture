COCKROACHDB_URL=cockroachdb://cockroach:@localhost:26257/userstore?sslmode=disable

sqlc:
	sqlc generate -f ./config/sqlc.yaml
run:
	go run cmd/main.go
migrate-up:
	migrate -database ${COCKROACHDB_URL} -path db/migration up
migrate-down:
	migrate -database ${COCKROACHDB_URL} -path db/migration down
godog:
	godog run test/features/register.feature
.PHONY: sqlc migrate-up migrate-down godog
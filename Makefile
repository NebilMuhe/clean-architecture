COCKROACHDB_URL=cockroachdb://cockroach:@localhost:26257/userstore?sslmode=disable

sqlc:
	sqlc generate -f ./config/sqlc.yaml
run:
	go run cmd/main.go
migrate-up:
	migrate -database ${COCKROACHDB_URL} -path db/migration up
migrate-down:
	migrate -database ${COCKROACHDB_URL} -path db/migration down
godog-run:
	godog run godog/features/register.feature
godog-test:
	godog test godog/tests/register_test.go
.PHONY: sqlc migrate-up migrate-down godog-run godog-test
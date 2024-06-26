COCKROACHDB_URL=cockroachdb://root:@localhost:26257/userstore?sslmode=disable

sqlc:
	sqlc generate -f ./config/sqlc.yaml
run:
	go run cmd/main.go
migrate-up:
	migrate -database ${COCKROACHDB_URL} -path internal/constants/db/migration up
migrate-down:
	migrate -database ${COCKROACHDB_URL} -path internal/constants/db/migration down
godog-run:
	godog run godog/features/
godog-test:
	godog test godog/tests/register_test.go
test:
	go test -v ./godog/tests
.PHONY: sqlc migrate-up migrate-down godog-run godog-test
.PHONY:
.SILENT:
.DEFAULT_GOAL := run

test:
	go test -coverprofile=cover.out -v ./...
	make test.coverage
test.coverage:
	go tool cover -func=cover.out | grep "total"
build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go
build_job:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/cron_job ./cmd/cron_job/main.go
run:
	go run ./cmd/app/main.go

run_cron:
	go run ./cmd/cron_job/main.go
swag:
	swag init --parseDependency --parseInternal --parseDepth 1 -g internal/app/app.go
gen:
	mockgen -source=internal/services/services.go -destination=internal/services/mocks/mock.go
	mockgen -source=pkg/googlesheets/sheets.go -destination=pkg/googlesheets/mocks/mock.go
	mockgen -source=pkg/mail/mail.go -destination=pkg/mail/mocks/mock.go
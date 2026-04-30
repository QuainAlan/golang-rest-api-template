setup:
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./cmd/server/main.go -o ./docs
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files

build-docker:
	docker compose build --no-cache

# Run API against local Docker DBs. Requires `.env` with JWT_SECRET_KEY and API_SECRET_KEY (see .env.example).
# Optional: set GIN_MODE=release in `.env` for the same security middleware as Docker Compose.
run-local:
	docker start dockerPostgres
	docker start dockerRedis
	docker start dockerMongo
	test -f .env || { echo >&2 "Missing .env — copy .env.example to .env and set secrets (>=32 bytes each)."; exit 1; }
	set -euo pipefail; \
	set -a && . ./.env && set +a; \
	export REDIS_HOST=localhost POSTGRES_HOST=localhost \
		POSTGRES_DB=go_app_dev POSTGRES_USER=docker POSTGRES_PASSWORD=password POSTGRES_PORT=5435; \
	go run cmd/server/main.go

up:
	docker compose up

down:
	docker compose down

restart:
	docker compose restart

build:
	go build -v ./...

test:
	go test -v ./... -race -cover

test-cover:
	PKGS=$$(go list ./... | grep -vE '(^|/)(cmd/server|docs|scripts)$$'); \
	go test -race -coverprofile=coverage.out $$PKGS; \
	go tool cover -html=coverage.out -o coverage.html

clean:
	docker stop go-rest-api-template
	docker stop dockerPostgres
	docker rm go-rest-api-template
	docker rm dockerPostgres
	docker rm dockerRedis
	docker image rm golang-rest-api-template-backend
	rm -rf .dbdata

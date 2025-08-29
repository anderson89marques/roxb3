GO = go

export_envs:
	export $(cat .env | xargs)

db_login:
	psql ${DATABASE_URL}	

# RUN LOCAL
.PHONY: build-api  run-api-dev

build-api:
	$(GO) build -o /cmd/api/api ./cmd/api/main.go

run-api-dev:
	$(GO) run cmd/api/main.go

build-ingest:
	$(GO) build -o /cmd/cli/cli ./cmd/cli/main.go

run-ingest-dev:
	$(GO) run cmd/cli/main.go

.PHONY: run-ingest  run-api
# Stock insgestion
run-ingest:
	docker-compose run --build --remove-orphans cli

# Run api
run-api:
	docker-compose up --build --force-recreate api -d

# Run db
run-db:
	docker-compose up --build db -d

# Database migrations
migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate:
	migrate -database ${DATABASE_URL} -path migrations up

# Swagger

swag:
	swag init ./cmd/api/main.go -o docs --parseDependency --parseInternal


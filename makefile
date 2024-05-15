# Get file .env
include .env
export $(shell sed 's/=.*//' .env)


DB_URL=postgresql://$$POSTGRES_USER:$$POSTGRES_PASSWORD@localhost:$$POSTGRES_PORT_MAPPING/$$POSTGRES_DB?sslmode=disable

# Folder constants
DOCKER_COMPOSE := docker-compose.yml

################# TEST #################
default:
	echo "$(DB_URL)"

################# DOCKER #################

run-build:
	docker-compose -f $(DOCKER_COMPOSE) up -d --build
	
run-down:
	docker-compose -f $(DOCKER_COMPOSE) down

################# SQLC #################
sqlc:
	sqlc generate

################# GO #################
run:
	go run ./cmd/main.go
test:
	go test -v -cover ./...


################# SQL #################
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
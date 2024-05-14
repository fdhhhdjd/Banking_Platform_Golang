# Get file .env
include .env
export $(shell sed 's/=.*//' .env)

# Folder constants
DOCKER_COMPOSE := docker-compose.yml

################# DOCKER #################
default:
	docker ps

run-build:
	docker-compose -f $(DOCKER_COMPOSE) up -d --build
	
run-down:
	docker-compose -f $(DOCKER_COMPOSE) down

################# SQLC #################
sqlc:
	sqlc generate

################# GO #################
test:
	go test -v -cover ./...



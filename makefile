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

dev:
	go run main.go

test:
	go test -v -cover ./...

################# SQL #################
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up


################# JWT #################
sereckey: 
	node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"

################# gRPC #################
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

proto-gateway:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--experimental_allow_proto3_optional \
	proto/*.proto

proto-swagger:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=banking \
	--experimental_allow_proto3_optional \
	proto/*.proto

evans:
	evans -r repl --host localhost --port 5005
	
.PHONY: proto evans proto-gateway

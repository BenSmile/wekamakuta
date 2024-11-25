# include .env

DB_URL := postgres://root:secret@localhost:5432/simple_bank?sslmode=disable
POSTGRES_USER := root
POSTGRES_PASSWORD := secret
POSTGRES_DB := simple_bank
.PHONY: postgres
postgres:
	@docker run --name simple-bank-pgdb -p 5432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:12-alpine
	# @docker run --network bank-network --name simple-bank-pgdb -p 5432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:12-alpine

.PHONY: createdb
createdb:
	@docker exec -it simple-bank-pgdb createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DB}

.PHONY: dropdb
dropdb:
	@docker exec -it simple-bank-pgdb dropdb ${POSTGRES_DB}

.PHONY: migrate-up
migrate-up:
	@migrate -path db/migration -database "${DB_URL}" -verbose up

.PHONY: migrate-up1
migrate-up1:
	@migrate -path db/migration -database "${DB_URL}" -verbose up 1

.PHONY: migrate-down
migrate-down:
	@migrate -path db/migration -database "${DB_URL}" -verbose down

.PHONY: migrate-down1
migrate-down1:
	@migrate -path db/migration -database "${DB_URL}" -verbose down 1

.PHONY: sqlc
sqlc:
	@sqlc generate

.PHONY: test
test:
	@go test -v -cover -count=1 ./...

.PHONY: dburl
dburl:
	@echo "Database URL: ${DB_URL}"

.PHONY: server
server:
	@go run main.go

.PHONY: mock
mock:
	@mockgen -package mockdb -destination db/mock/store.go github.com/bensmile/wekamakuta/db/sqlc Store

.PHONY: db_docs
db_docs:
	@dbdocs build doc/db.dbml

.PHONY: db_schema
db_schema:
	@dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PHONY: proto
proto:
	@rm -f pb/*go
	@rm -f doc/swagger/*.swagger.json
	@protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
   		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simplebank\
    	proto/*.proto
	@statik -src=./doc/swagger -dest=./src

.PHONY: evans
evans:
	@evans --host localhost --port 9090 -r repl

include .env

.PHONY: postgres
postgres:
	# @docker run --name simple-bank-pgdb -p  "5432:5432"  -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -e POSTGRES_DB=${POSTGRES_DB} -d postgres:12-alpine 
	@docker run --name simple-bank-pgdb -p  "5432:5432"  -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:12-alpine 

.PHONY: createdb
createdb:
	@docker exec -it simple-bank-pgdb createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DB}

.PHONY: dropdb
dropdb:
	@docker exec -it simple-bank-pgdb dropdb ${POSTGRES_DB}


.PHONY: migrate-up
migrate-up:
	@migrate -path db/migration -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	@migrate -path db/migration -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable" -verbose down

.PHONY: sqlc
sqlc:
	@sqlc generate

.PHONY: test
test:
	@go test -v -cover ./...
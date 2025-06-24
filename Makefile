CURRENT_DIR=$(shell pwd)

-include .env

PSQL_CONTAINER_NAME?=postgresql
PROJECT_NAME?=menu-mono
PSQL_URI=postgres://postgres:postgres@localhost:5444/${PROJECT_NAME}?sslmode=disable

createdb:
	docker exec -it ${PSQL_CONTAINER_NAME} createdb --username=postgres --owner=postgres ${PROJECT_NAME}

dropdb:
	docker exec -it ${PSQL_CONTAINER_NAME} dropdb -U postgres ${PROJECT_NAME}

# make new_migration name=init_schema
new_migration:
	goose -dir db/migrations create -s $(name) sql 

migrate_up:
	goose -dir db/migrations postgres "${PSQL_URI}" up

migrate_down:
	goose -dir db/migrations postgres "${PSQL_URI}" down

migrate_status:
	goose -dir db/migrations postgres "${PSQL_URI}" status

sqlc:
	sqlc generate
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

proto:
	rm -f generated/**/*.go
	rm -f doc/swagger/*.swagger.json
	mkdir -p generated
	protoc \
		--proto_path=protos --go_out=generated --go_opt=paths=source_relative \
		--go-grpc_out=generated --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=generated --grpc-gateway_opt=paths=source_relative,allow_delete_body=true \
		--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=swagger_docs,use_allof_for_refs=true,disable_service_tags=false,allow_delete_body=true \
		--validate_out="lang=go,paths=source_relative:generated" \
	protos/**/*.proto
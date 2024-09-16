DB_CONTAINER=postgres
DB_USER=root
DB_PASSWORD=root123
DB_IP_ADDRESS=localhost
DB_PORT=5432
DB_NAME=simple_bank
DB_TERMINATE_COMMAND=psql --username=$(DB_USER) --dbname=postgres --command="SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname='$(DB_NAME)';"
DB_CREATE_COMMAND=createdb --username=$(DB_USER) $(DB_NAME)
DB_DROP_COMMAND=dropdb --username=$(DB_USER) --if-exists $(DB_NAME)
IMAGE_NAME=postgres
IMAGE_TAG=alpine3.19

# This target is used to start a PostgreSQL server in a Docker container.
# Make sure you have Docker installed on your machine and postgres image with this this tag.
postgres-start:
	docker run --name $(DB_CONTAINER) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -p $(DB_PORT):5432 -d $(IMAGE_NAME):$(IMAGE_TAG)

postgres-stop:
	docker stop $(DB_CONTAINER) || true
	docker rm $(DB_CONTAINER) || true
# This target is used to create a new database if database does not exist already.
createdb:
	docker exec -it $(DB_CONTAINER) $(DB_CREATE_COMMAND)

# This target will connect to the PostgreSQL server and drop the specified database.
# Usage: make dropdb
# Ensure that the database server is running and you have the necessary permissions to drop the database.
dropdb:
	docker exec -it $(DB_CONTAINER) $(DB_TERMINATE_COMMAND)
	docker exec -it $(DB_CONTAINER) $(DB_DROP_COMMAND)

migrateup:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_IP_ADDRESS):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_IP_ADDRESS):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down


.PHONY: createdb dropdb postgres-start postgres-stop migrateup migratedown


include .env
export

export PROJECT_ROOT=${shell pwd}
export LOGGER_FOLDER=${PROJECT_ROOT}/log

env-up:
	docker compose up -d --build && \
	go run cmd/main.go


migrate-up:
	@migrate \
	-path ${PROJECT_ROOT}/migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable	\
	up

migrate-down:
	@migrate \
	-path ${PROJECT_ROOT}/migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable	\
	down


# Makefile

.PHONY: compose-up compose-down compose-loader

## Bring up Mongo, app and loader in background
compose-up:
	docker-compose up --build

## Tear down Docker Compose stack
compose-down:
	docker-compose down

## Run loader via Docker Compose
compose-loader:
	docker-compose up --build loader

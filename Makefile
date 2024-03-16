SHELL := /bin/bash

DC_FILE="docker-compose.yml"
DC_RESOURCE_DIR=".compose"

test:
	go test -v ./...

prepare-compose:
	test -d $(DC_RESOURCE_DIR) || mkdir $(DC_RESOURCE_DIR) || true
	test -f $(DC_RESOURCE_DIR)/config.yml || cp config.example.yml $(DC_RESOURCE_DIR)/config.yml || true

up: prepare-compose
	docker-compose -f $(DC_FILE) up -d --build

down:
	docker-compose down

integration-test: up
	go test -tags integration -v ./internal/it
	docker-compose down --volumes


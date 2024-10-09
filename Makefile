#!/usr/bin/make

SHELL = /bin/sh

init:
	swag init
	docker compose up -d postgres redis
	docker compose exec -it postgres psql -U postgres -c "CREATE DATABASE db;"
	docker compose up -d --build app

test:
	go test -v ./...

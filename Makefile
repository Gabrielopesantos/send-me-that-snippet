.PHONY: local run build test

local:
	echo "Starting local environment"
	docker-compose -f docker-compose.local.yml up -d

local-down:
	echo "Deleting local environment"
	docker-compose -f docker-compose.local.yml down -v --remove-orphans

run:
	go run ./cmd/server/main.go

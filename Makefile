run:
	DNV=dev go run cmd/app/main.go

.PHONY: dev-up dev-down dev-logs
docker-up:
	docker compose -f docker-compose.yml up -d
docker-down:
	docker compose -f docker-compose.dev.yml down
docker-logs:
	docker compose -f docker-compose.dev.yml logs -f
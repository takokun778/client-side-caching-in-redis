.PHONY: up
up:
	@docker compose --project-name client-side-caching-in-redis --file ./.docker/compose.yaml up -d

.PHONY: down
down:
	@docker compose --project-name client-side-caching-in-redis down

.PHONY: cli
cli:
	@docker exec -it client-side-caching-in-redis redis-cli

.PHONY: a
a:
	@go run cmd/a/main.go

.PHONY: b
b:
	@go run cmd/b/main.go

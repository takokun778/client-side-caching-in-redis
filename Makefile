.PHONY: up
up:
	@docker compose --project-name client-side-caching-in-redis --file ./.docker/compose.yaml up -d

.PHONY: down
down:
	@docker compose --project-name client-side-caching-in-redis down

.PHONY: cli
cli:
	@docker exec -it redis redis-cli

.PHONY: a
a:
	@curl http://localhost:8081/set?key=key\&val=a
	@curl http://localhost:8081/get?key=key

.PHONY: b
b:
	@curl http://localhost:8082/get?key=key

.PHONY: test
test:
	@go test -v ./test/... -count=1

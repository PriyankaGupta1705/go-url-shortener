.PHONY: local build down logs test clean psql redis-cli

local:
	docker compose -f docker-compose.yml -f docker-compose.local.yml up --build

down:
	docker compose down

logs:
	docker compose logs -f app

test:
	go test ./...

clean:
	docker compose down -v
	docker system prune -f

# open postgres shell
psql:
	docker compose exec database psql -U admin -d urlshortener

# open redis cli
redis-cli:
	docker compose exec redis redis-cli
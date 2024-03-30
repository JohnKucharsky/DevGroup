dev:
	fresh
dev-db:
	docker compose -f compose-dev.yaml up -d
prod:
	docker compose up --build
migrate:
	cd db/migrations; goose postgres postgres://postgres:pass@localhost:5432/data up
migrate-down:
	cd db/migrations; goose postgres postgres://postgres:pass@localhost:5432/data down
api:
	cd ./public; redocly build-docs ./openapi.yaml --output=index.html
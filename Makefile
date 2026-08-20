.PHONY: db-up db-down migrate-up migrate-down backend-run backend-test frontend-dev frontend-test check

db-up:
	docker compose up -d mysql

db-down:
	docker compose down

migrate-up:
	cd backend && go run -tags mysql github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3 -path migrations -database "mysql://article_user:article_password@tcp(localhost:3306)/article" up

migrate-down:
	cd backend && go run -tags mysql github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3 -path migrations -database "mysql://article_user:article_password@tcp(localhost:3306)/article" down 1

backend-run:
	cd backend && go run ./cmd/api

backend-test:
	cd backend && go test ./...

frontend-dev:
	cd frontend && npm run dev

frontend-test:
	cd frontend && npm run test

check:
	cd backend && go vet ./... && go test ./... && go build ./cmd/api
	cd frontend && npm run lint && npm run typecheck && npm run test -- --run && npm run build


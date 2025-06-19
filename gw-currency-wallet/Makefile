MIGRATIONS_DIR=./migrations

DSN=postgres://user:123456@localhost:5432/postgres?sslmode=disable

migrate-create-%:
	goose -dir $(MIGRATIONS_DIR) create $(subst migrate-create-,,$@) sql

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DSN)" down

migrate-reset:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DSN)" reset

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DSN)" status


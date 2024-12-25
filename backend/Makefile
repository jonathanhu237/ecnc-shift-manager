DSN = postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable

new_migration:
	migrate create -seq -ext=.sql -dir=./migrations ${name}

migrate_up:
	migrate -path ./migrations -database ${DSN} up

migrate_down:
	migrate -path ./migrations -database ${DSN} down ${step}

migration_version:
	migrate -path ./migrations -database ${DSN} version

set_migration_version:
	migrate -path ./migrations -database ${DSN} force ${version}
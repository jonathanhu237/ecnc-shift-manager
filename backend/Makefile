psql:
	psql ${ECNC_SHIFT_MANAGER_DB_DSN}

new_migration:
	migrate create -seq -ext=.sql -dir=./migrations ${name}

migrate_up:
	migrate -path ./migrations -database ${ECNC_SHIFT_MANAGER_DB_DSN} up

migrate_down:
	migrate -path ./migrations -database ${ECNC_SHIFT_MANAGER_DB_DSN} down ${step}
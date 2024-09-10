.PHONY: watch run db kill up down force air

run:
	go run main.go

watch:
	reflex -r '\.go$$' -s -- sh -c 'make run'

air:
	air

kill:
	npx kill-port 8080

up:
	go run db/migrate/up.go

down:
	go run db/migrate/down.go

force:
	go run db/migrate/force.go

db:
	@if [ -z "$(n)" ]; then \
            echo "Error: name is not set. Use 'make db n=yourfilename'"; \
            exit 1; \
	fi
	migrate create -ext sql -dir db/migration -seq $(n)

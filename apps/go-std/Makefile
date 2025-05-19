.PHONY: migrate-db dump-schema

# Path configurations
PRISMA_SCHEMA=db/prisma/schema.prisma
SQL_SCHEMA=db/schemas/schema.sql

# Migrate the database using Prisma
push-db:
	@echo "Running Prisma push..."
	prisma db push --schema=$(PRISMA_SCHEMA)

pull-db:
	@echo "Running Prisma pull..."
	prisma db pull --schema=$(PRISMA_SCHEMA)

# Generate schema diff from empty state
dump-schema:
	@echo "Generating schema diff to $(SQL_SCHEMA)..."
	prisma migrate diff \
		--from-empty \
		--to-schema-datamodel=$(PRISMA_SCHEMA) \
		--script > $(SQL_SCHEMA)

start-cloudflared:
	cloudflared access tcp --hostname pg.mlcr.us --url localhost:5433


# Run both migration and schema dump
db-update: push-db dump-schema


test-router:
	go test -bench=. ./routes

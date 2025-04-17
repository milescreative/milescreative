.PHONY: migrate-db dump-schema

# Path configurations
PRISMA_SCHEMA=db/prisma/schema.prisma
SQL_SCHEMA=db/schemas/schema.sql

# Migrate the database using Prisma
migrate-db:
	@echo "Running Prisma migrations..."
	prisma migrate dev --schema=$(PRISMA_SCHEMA)

# Generate schema diff from empty state
dump-schema:
	@echo "Generating schema diff to $(SQL_SCHEMA)..."
	prisma migrate diff \
		--from-empty \
		--to-schema-datamodel=$(PRISMA_SCHEMA) \
		--script > $(SQL_SCHEMA)

# Run both migration and schema dump
db-update: migrate-db dump-schema

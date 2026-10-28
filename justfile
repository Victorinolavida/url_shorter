db_url := env_var_or_default("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/url_shorten?sslmode=disable")
# Directory containing dbmate migration files.
migrations := "db/migrations"
pg_container := "url-shorten-postgres"

# Schema file consumed by sqlc, regenerated from the migrated database.
schema_file := "db/schema.sql"

# sqlc configuration file.
sqlc_config := "db/sqlc.yaml"


run:
    @go run ./cmd/api/...

run-worker:
    @go run ./cmd/worker/...

# Apply all pending DB migrations.
migrate: db-create
    dbmate --url "{{db_url}}" --migrations-dir "{{migrations}}" up

# Regenerate sqlc schema.sql from the migrated database.
schema: migrate
    docker exec {{pg_container}} pg_dump --schema-only --no-owner --no-privileges \
        --exclude-table=schema_migrations "{{db_url}}" | grep -v '^\\' > "{{schema_file}}"

# Generate sqlc code from queries and schema.
generate: schema
    sqlc generate --file "{{sqlc_config}}"

# Roll back the last migration batch.
down: db-create
    dbmate --url "{{db_url}}" --migrations-dir "{{migrations}}" down

# Drop the database and reapply all migrations.
clean: db-create
    dbmate --url "{{db_url}}" --migrations-dir "{{migrations}}" drop
    dbmate --url "{{db_url}}" --migrations-dir "{{migrations}}" up

# Create a new migration file.
new name:
    dbmate --url "{{db_url}}" --migrations-dir "{{migrations}}" new {{name}}

# Show pending/applied migrations.
status: db-create
    dbmate --url "{{db_url}}" --migrations-dir "{{migrations}}" status


# Create the database if it does not exist yet.
[private]
db-create:
    dbmate --url "{{db_url}}" create 2>/dev/null || true


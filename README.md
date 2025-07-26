# Gator

sql url - "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

psql -U postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

## rerun migrations

# First, roll back all migrations

goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" reset

# Then apply them again

goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up

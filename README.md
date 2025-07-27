# Gator

RSS feed aggregator CLI tool built with Go and PostgreSQL.

## Prerequisites

- Go 1.19+ installed
- PostgreSQL installed and running

## Installation

### Install from source using go install:

```bash
go install github.com/Alexeychuk/Gator@latest
```

This will install the `gator` binary to your `$GOPATH/bin` directory (usually `~/go/bin`).

### Make sure Go bin is in your PATH:

```bash
# Add this to your ~/.zshrc or ~/.bashrc
export PATH="$HOME/go/bin:$PATH"
```

Then reload your shell:

```bash
source ~/.zshrc
```

### Verify installation:

```bash
gator --help
```

## Usage

First, configure your database connection and register a user:

```bash
# Register a new user
gator register <username>

# Login as user
gator login <username>

# Add an RSS feed
gator addfeed "Feed Name" "https://example.com/rss.xml"

# Follow a feed
gator follow "https://example.com/rss.xml"

# Start aggregating feeds (every 60 seconds)
gator agg 60s

# Browse posts
gator browse 10
```

## Development

### Database Setup

sql url - "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

### Migrations

```bash
# First, roll back all migrations
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" reset

# Then apply them again
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

### Build from source

```bash
git clone https://github.com/Alexeychuk/Gator.git
cd Gator
go build -o gator
./gator
```

To run yo need Go and Postgres installed

dev notes

sql url - "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

psql -U postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"

## rerun migrations

# First, roll back all migrations

goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" reset

# Then apply them again

goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up

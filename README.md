# go-fullstack-boilerplate

A template to scaffold a fullstack golang web application.

## Features

-   Standard Go Project [Layout](https://github.com/golang-standards/project-layout)
-   Postgresql database using database/sql with [pgx](https://pkg.go.dev/github.com/jackc/pgx/stdlib) driver
-   [Router](https://github.com/ferdiebergado/go-express) based on net/http ServeMux
-   Templating based on html/template
-   Optimized css and js builds
-   Database migrations
-   Hot reloading

## Requirements

-   Go version 1.22 or higher
-   Docker or Podman
-   [esbuild](https://esbuild.github.io/getting-started/)

## Usage

1. Install the cli tools.

```sh
make install
```

2. Rename .env.example to .env.

```sh
mv .env.example .env
```

3. Change the database password (DB_PASS).

```.env
# .env
DB_PASS=CHANGE_ME
```

4. Start the database.

```sh
make db
```

5. Run the server in development mode.

```sh
make dev
```

6. Open the web application at [localhost:8888](http://locahost:8888).

## Migrations

### Creating Migrations

Run the migration target with the name argument set to the name of the migration.

```sh
make migration name=create_users_table
```

### Running Migrations

Run the migrate target.

```sh
make migrate
```

### Rolling Back Migrations

Run the rollback target.

```sh
make rollback
```

## Running Tests

```sh
make test
```

## Other Tasks

Consult the Makefile.

## Linting

This project comes with a golangci-lint config file. Just install golangci-lint to use it.

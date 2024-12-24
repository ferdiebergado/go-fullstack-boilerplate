# go-fullstack-boilerplate

[![Go Report Card](https://goreportcard.com/badge/github.com/ferdiebergado/go-fullstack-boilerplate)](https://goreportcard.com/report/github.com/ferdiebergado/go-fullstack-boilerplate)

A template to scaffold a fullstack go web application.

## Features

-   Standard Go Project [Layout](https://github.com/golang-standards/project-layout)
-   Postgresql database using database/sql with [pgx](https://pkg.go.dev/github.com/jackc/pgx/stdlib) driver
-   [Router](https://github.com/ferdiebergado/goexpress) based on net/http ServeMux
-   HTML templating using html/template
-   Typescript support out-of-the-box
-   Database migrations
-   Hot reloading during development
-   [nginx](https://nginx.org/en/) as web server and reverse proxy configured for high-performance

## Requirements

-   Go version 1.22 or higher
-   Docker or Podman

## Usage

### Step 1

Rename .env.example to .env.

```sh
mv .env.example .env
```

### Step 2

Change the database password (DB_PASSWORD).

```.env
# .env
DB_PASSWORD=CHANGE_ME
```

Optionally, you can also set the user and the database.

### Step 3

Deploy the application.

```sh
make
```

### Step 4

Browse the application at [localhost:8080](http://locahost:8080).

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

## Bundling assets

### Bundle for development

```sh
make bundle
```

### Watch mode for css files

```sh
make watch-css
```

### Watch mode for typescript/javascript files

```sh
make watch-ts
```

### Bundle for production

```sh
make bundle-prod
```

## Running Tests

```sh
make test
```

## Other Tasks

### Interact with the database using psql

```sh
make psql
```

### Restart a service

Provide a service argument to the restart target.

-   Restart the reverse proxy:

```sh
make restart service=proxy
```

### Stop all the running containers

```sh
make stop
```

## Linting

This project comes with a golangci-lint config file. Just install golangci-lint and enable it as the default linter on your editor of choice.

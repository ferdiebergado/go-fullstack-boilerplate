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

Deploy the application.

```sh
make
```

### Step 3

Browse the application at [localhost:8080](http://locahost:8080).

## Migrations

### Create Migrations

Run the migration target with the name argument set to the name of the migration.

```sh
make migration name=create_users_table
```

### Run Migrations

```sh
make migrate
```

### Rollback Migrations

```sh
make rollback
```

### Recover from a Failed Migration

When a migration fails, fix the error and force the version of the failed migration.
Then run the migration again.

```sh
make force version=1
make migrate
```

## Bundling of Assets

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

## Tests

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

```sh
make restart service=proxy
```

### Stop all the running containers

```sh
make stop
```

## Linting

This project comes with a golangci-lint config file. Just install golangci-lint and enable it as the default linter on your editor of choice.

## License

This project is distributed under the MIT License. See [LICENSE](https://github.com/ferdiebergado/go-fullstack-boilerplate/blob/main/LICENSE) for more details.

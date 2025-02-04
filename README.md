# go-fullstack-boilerplate

![Github Actions](https://github.com/ferdiebergado/go-fullstack-boilerplate/actions/workflows/go.yml/badge.svg?event=push) ![Github Actions](https://github.com/ferdiebergado/go-fullstack-boilerplate/actions/workflows/security.yml/badge.svg?event=push) [![Go Report Card](https://goreportcard.com/badge/github.com/ferdiebergado/go-fullstack-boilerplate)](https://goreportcard.com/report/github.com/ferdiebergado/go-fullstack-boilerplate)

A template to scaffold a fullstack go web application.

## Features

-   Standard Go Project [Layout](https://github.com/golang-standards/project-layout)
-   Postgresql database using database/sql with [pgx](https://pkg.go.dev/github.com/jackc/pgx/stdlib) driver
-   [Router](https://github.com/ferdiebergado/goexpress) based on net/http ServeMux
-   HTML templating using html/template
-   Typescript support out-of-the-box
-   [Toolkit](https://github.com/ferdiebergado/gopherkit) that makes common tasks easier
-   Database migrations
-   Hot reloading during development
-   [nginx](https://nginx.org/en/) as web server and reverse proxy configured for high-performance
-   Docker deployment

## Requirements

-   Go version 1.22 or higher
-   Docker

## Getting Started

### Step 1

Rename .env.example to .env.

```sh
mv .env.example .env
```

### Step 2

Deploy the application.

```sh
make dev
```

### Step 3

Browse the application at [localhost:8080](http://localhost:8080).

## Migrations

### Create Migrations

Run the migration target with the name of the migration as argument.

```sh
make migration create_users_table
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
make force 1
make migrate
```

## Bundling Assets

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

Run unit tests.

```sh
make test
```

Run integration tests.

```sh
make integration
```

## Other Tasks

### Interact with the database using psql

```sh
make psql
```

### Restart a service

Provide a service as argument to the restart target.

```sh
make restart proxy
```

### Stop all the running containers

```sh
make stop
```

### Help

View the usage information by running make.

```sh
make
```

## TODOs

-   [x] Health endpoint
-   [x] Login with email and password
-   [ ] Email verification
-   [ ] Secure Cookie Session Management
-   [ ] Login with Google (OAuth2)
-   [ ] Authorization
-   [ ] Audit logs
-   [ ] Database query caching
-   [ ] Environment Page (go version, drivers, env, os kernel, etc.)
-   [ ] Cache busting for assets

## Linting

This project comes with a golangci-lint config file. Just install golangci-lint and enable it as the default linter on your editor of choice.

## License

This project is distributed under the MIT License. See [LICENSE](https://github.com/ferdiebergado/go-fullstack-boilerplate/blob/main/LICENSE) for more details.

## Tech Stack

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![TypeScript](https://img.shields.io/badge/typescript-%23007ACC.svg?style=for-the-badge&logo=typescript&logoColor=white) ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) ![Nginx](https://img.shields.io/badge/nginx-%23009639.svg?style=for-the-badge&logo=nginx&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

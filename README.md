# YouChoose

A http server written in Go for the YouChoose mobile app.

## Development

### Environment

We utilise [development containers](https://containers.dev) to provide a fully featured, reproducible development environment that is easily shareable.

The config for our container is located in `.devcontainer` and comprises of a `docker-compose` and `Dockerfile` which sets up the environment.

To connect to the environment in VSCode, open the command palette (`CMD-SHIFT-P` | `CTRL-SHIFT-P`) and type `Dev Containers: Rebuild and Reopen in Container`. The development environment will be created, it may take some time to complete.

### Makefile

We utilise `make` to store useful commands for building, running and testing our app.

To see the list of availible actions:

```bash
development@94242c61cf2b:/app$ # inside your container...

# Print all available make targets
make help
```

### Postgres

A PostgreSQL database is automatically started and exposed on `localhost:5432`. You can connect to it with either a database client on your host machine or using `psql` within the development container. The configuration for the postgres instance is location inside the `.devcontainer/docker-compose.yml` file.

### sqlc

[sqlc](https://github.com/sqlc-dev/sqlc) is installed in the container and is used to generate type-safe code from raw SQL. Our sql files are split into two; queries and migrations. Queries are SQL queries with special comments that sqlc can read and use to generate our database access functions. Migrations contains our up and down migration scripts, which will be handled by tern. The configuration for sqlc is contained in the `sqlc.yml` file.

### Tern

[Tern](https://github.com/jackc/tern) is a migration tool for PostgreSQL. The configuration for tern is contained in the `tern.conf` file.

### IntegreSQL

A [IntegreSQL](https://github.com/allaboutapps/integresql) service is automatically started in the background. It is used to manage isolated PostgreSQL databses for our integration tests. The configuration for the integres instance is location inside the `.devcontainer/docker-compose.yml` file.

### gotestsum

[gotestsum](https://github.com/gotestyourself/gotestsum) is installed in the container as our main test runner. It is optimized for output test results in a human-readable way.

~~~bash
development@94242c61cf2b:/app$ # inside your container...

# run gotestsum and display the test results for each package
make go-test-by-pkg
~~~

### golangci-lint

[golangci-lint](https://github.com/golangci/golangci-lint) is installed in our container and is used for linting purposes. It comes with many additional linters but we are using just a small subset. The configuration for golangci-lint is location in the `.golangci.yml` file.

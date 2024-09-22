# YouChoose

A http server written in Go for the YouChoose mobile app.

## Development

### Environment

We utilise [development containers](https://containers.dev) to provide a fully featured, reproduciple development environment that is easily shareable.

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

A PostgreSQL database is automatically started and exposed on `localhost:5432`. You can connect to it with either a database client on your host machine or using `psql` within the development container.

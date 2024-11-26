# You Choose

## Development

### Local Setup

This project utilises [devenv](https://devenv.sh) to create a declarative, reproducible development environment. Follow the [getting started guid](https://devenv.sh/getting-started/) in order to install devenv on your system. We also use [direnv](https://direnv.net) to automatically activate our devenv shell.

If you do not wish to use our devenv environment, then ensure you have all of the dependencies specified inside the `devenv.nix` file.

### Commands

#### fhrs

To populate our restaurants table, we fetch the data from the free and open data set provided by the Food Standard Agency. The data is updated regularly and uploaded in xml format to [here](https://ratings.food.gov.uk/open-data). We can fetch the data and upload it to our database with the following command:

~~~shell
$ go run cmd/fhrs/main.go --database-uri="${DATABASE_URL}"
~~~

which will fetch the data for Glasgow City region. To specify a different region, pass the `--region-id` argument.

### Tools

#### Database Migrations - Tern

We use [Tern](https://github.com/jackc/tern) to handle database migrations. The migrations and our tern config can be found in `backend/migrations/`. Our `.env` file exports two tern environment variables:

~~~shell
TERN_CONFIG=path-to-project/backend/migrations/tern.conf
TERN_MIGRATIONS=path-to-project/backend/migrations/
~~~

which allows us to run tern commands from anywhere in the command line.

#### Database - PostgreSQL

We use [PostgreSQL 16](https://www.postgresql.org/docs/16/index.html) as our database. With Tern and the excellent driver for Go, [pgx](https://github.com/jackc/pgx) we can do just about anything.

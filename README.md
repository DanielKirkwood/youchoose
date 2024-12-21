# You Choose

## Development

### Local Setup



### Commands

#### fhrs

To populate our restaurants table, we fetch the data from the free and open data set provided by the Food Standard Agency. The data is updated regularly and uploaded in xml format to [here](https://ratings.food.gov.uk/open-data). We can fetch the data and upload it to our database with the following command:

~~~shell
$ go run cmd/fhrs/main.go --database-uri="${DATABASE_CONNSTRING}"
~~~

which will fetch the data for Glasgow City region. To specify a different region, pass the `--region-id` argument.

### Tools

#### Database Migrations - Tern

We use [Tern](https://github.com/jackc/tern) to handle database migrations. The migrations and our tern config can be found in `backend/migrations/`. Our `.envrc` file exports two tern environment variables:

~~~shell
export TERN_CONFIG=path-to-project/backend/migrations/tern.conf
export TERN_MIGRATIONS=path-to-project/backend/migrations/
~~~

which allows us to run tern commands from anywhere in the command line.

#### Database - PostgreSQL

We use [PostgreSQL 16](https://www.postgresql.org/docs/16/index.html) as our database. With Tern and the excellent driver for Go, [pgx](https://github.com/jackc/pgx) we can do just about anything.

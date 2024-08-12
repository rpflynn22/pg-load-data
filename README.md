# pg-load-data


pg-load-data is a simple project that helps to run a local instance of Postgres
in Docker and uses a schema file and a go program in concert to bootstrap the
database and load data.

I used this to debug read queries that only show problems with a large amount of
data present.

## Usage

- Create a new folder under `./cases`
- Add a `schema.sql` file with sql commands to setup the db however needed
- Add a `main.go` file that uses `pkg/database` to initialize a db connection
and `pkg/insertdriver` to run the data load
  - See existing folders in `./cases` for examples
- Export your new directory name like `export CASE=<dirname>` and run make commands
to load the data: `make pg-up apply-schema insert`

# Testing pgbouncer connection management (its stability)

The following instructions apply for Ubuntu (for other OS you might use similar approach).

## Dependencies

### Postgres
You'll need to have Postgres DB installed (and working on its default 5432 port). 

You'll need to create Postgres user (e.g. as described [here](https://support.rackspace.com/how-to/postgresql-creating-and-dropping-roles/)).

You'll need to create Postgres DB.
You can use the following script to create Postgres DB (named `test_db`, will be used later for everything):
```bash
echo "DROP DATABASE IF EXISTS test; CREATE DATABASE test $(COLLATE);" | PGPASSWORD=test psql postgres -U test -h 127.0.0.1 -w
echo "DROP SCHEMA IF EXISTS test cascade; CREATE SCHEMA test AUTHORIZATION test;" \
"ALTER ROLE test SET search_path TO test;" \
"GRANT USAGE ON SCHEMA test TO test;" \
"ALTER DATABASE test SET timezone TO 'UTC';" \
"REVOKE ALL ON SCHEMA public FROM public;" | PGPASSWORD=test psql test -U test -h 127.0.0.1 -w
```

### pgbouncer

You'll need to install pgbouncer. On Ubuntu you can do:
```bash
sudo apt-get install pgbouncer
```
Required configuration files already exist in this repository.

### Go

The following libraries are used by Go test script:
```bash
go get github.com/jmoiron/sqlx
go get github.com/stretchr/testify
go get github.com/LasTshaMAN/Go-Execute/jobs
```

## Running tests

Run pgbouncer with:
```bash
pgbouncer pg_b_cfg.ini
```
You can connect to Postgres through pgbouncer via console:
```bash
PGPASSWORD=test psql -h 127.0.0.1 -p 6543 -U test test_db
```
Or you can run Go stress test:
```bash
go test test_test.go -v
```

// TODO
// About network partitioning
## Enable network partitioning
You might want to simulate real-world network partitioning by configuring it like this:
```bash

```

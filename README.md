# Transfer API
This service act as a money transfer service that serve user transfer/disbursement to a bank account.

## Dependencies
- Go 1.20
- PostgreSQL 16.2
- Docker (optional)
- [golang-migrate](https://github.com/golang-migrate/migrate)

## Setup
1. Clone this repository
```bash
git clone git@github.com:rakhmatullahyoga/transfer.git
```
2. Setup environment variables. You can easily setup environment variables using `.env` file by copying from the `env.sample` file and modifying the `.env` file
```bash
make env
```
3. Setup dependencies via Docker (optional)
```bash
docker-compose up -d
```
4. Run database migration to setup database schema. To run the migration, please follow the instruction on the [developer page](https://pkg.go.dev/github.com/golang-migrate/migrate/cli#section-readme)
```bash
make bin/migrate
./bin/migrate [migrate commands]
```
5. Compile the project
```bash
make compile
```
6. Run the application
```bash
make run
```

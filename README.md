# Transfer API
This service act as a money transfer service that serve user transfer/disbursement to a bank account.

## Dependencies
- Go 1.20
- PostgreSQL 16.2
- Docker (optional)
- [golang-migrate](https://github.com/golang-migrate/migrate)

## Project structure

```
root
|- bin # binary files, contains application binary and additional tools such as `migrate`
|- db # db related scripts
|  |- migrations # migration sql files
|     |- x.down.sql
|     |- x.up.sql
|- transfer # go package for transfer API core business process: account inquiry, transfer process, and callback endpoint
|  |- handler # http and other delivery layer handler
|     |- xxx.go
|  |- repository # data repository layer functions, each data source have separate files
|     |- xxx.go
|  |- usecase # usecase layer functions
|     |- xxx.go
|  |- domain.go # domain's model and architecture layer definitions
|- docker-compose.yml # system dependency
|- main.go # application entrypoint
```

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

## Mock API
This service used [mockAPI](https://mockapi.io) to simulate external dependency to bank service. We use `accounts` and `transfers` as resources in the mockAPI. The base URL for this project is `https://65d4055f522627d50109c299.mockapi.io/api/v1`

### Account validation
Getting account information of an account number, including it's owner name.

<details>
 <summary><code>GET</code> <code><b>/accounts</b></code> <code>(get accounts associated with an account number)</code></summary>

##### Parameters

> | name          | type     | data type | description           |
> |---------------|----------|-----------|-----------------------|
> | accountNumber | optional | string    | A bank account number |


##### Responses

> | http code | content-type       | response                                                                              |
> |-----------|--------------------|---------------------------------------------------------------------------------------|
> | `200`     | `application/json` | `[{"owner":"Dwayne Swift V","accountNumber":"37976402","balance":"673.51","id":"1"}]` |

</details>

### Transfer money
Add transfer data to mockAPI database.

<details>
 <summary><code>POST</code> <code><b>/transfers</b></code> <code>(add transfer data to mockAPI database)</code></summary>

##### Request body (JSON)

> | name      | type     | data type | description                     |
> |-----------|----------|-----------|---------------------------------|
> | accountId | required | string    | Transfer destination account ID |
> | amount    | required | string    | Transfer amount                 |


##### Responses

> | http code | content-type       | response                                                                                               |
> |-----------|--------------------|--------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | `{"createdAt":"2024-02-19T11:00:34.135Z","amount":"100.00","success":false,"id":"26","accountId":"2"}` |

</details>

## Test the application
You can test the application by making http request to transfer service as described in the attached [Postman collection](Transfer-API.postman_collection.json).

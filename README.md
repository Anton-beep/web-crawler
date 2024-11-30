# Production

To run a project in production mode

- Set up environment variables in the `configs/Docker.env` file use `configs/Docker.env.template` for this
- Install [Docker](https://www.docker.com)
- Run the command
```bash
docker compose -f .\deployments\docker-compose.yml up
```

# Dev

## Pre-setting

To run a project in development mode
- Set up environment variables in the `configs/.env` use `configs/.env.template` for this.
- Install [Docker](https://www.docker.com)
- Run the command
```bash
docker compose -f .\deployments\docker-compose-dev.yml up
```
In this mode, all third-party containers will be deployed. The ports will be forwarded to localhost.
## How to start services

- Place the executable file at `cmd/<your_service>/<file_name>.go`
- Follow this path
```bash
cd cmd/<your_service>
```
- Run the executable file
```bash
go run <file_name>.go
```

## How to write tests

If you are writing unit tests, then place them along the path `test/<name_of_package>/<test_name>_test.go`

If you are writing integration tests, then place them in the path `test/integrationTests/<test_name>_test.go`.
At the beginning of each test, check that the environment variable `IS_RAN_IN_DOCKER` is true. 
This ensures that tests run alongside running containers.

For example, you can do this as follows

```go

package integrationTests_test

import (
	"testing"
	"web-crauler/internal/config"
	"web-crauler/internal/repository"
)

var (
	db  models.DataBase
	cfg *config.Config
)

func TestMain(m *testing.M) {
	db = repository.NewDB()
	cfg = config.NewConfig()
	code := m.Run()
	os.Exit(code)
}

func SomeTests(t *testing.T) {
	if !cfg.IsRanInDocker {
		return
	}
	// ...
}
```

## Folder structure

Follow this [standard structure](https://github.com/golang-standards/project-layout) of a go project


## Linter

```shell
golangci-lint run -c configs/.golangci.yml
```


## Testing

### Unit tests

To run unit tests, run the command
```bash
go test -cover -v ./...
```
The console will display the testing results, as well as the percentage of test coverage

### Integration tests

To run integration tests,
- Set up environment variables in the `configs/.env`, use `configs/.env.template` for this
- Set the value of the `RUN_INTEGRATION_TESTS` variable to `true`
- Install [Docker](https://www.docker.com)
- Run the command
```bash
docker compose -f .\deployments\docker-compose-dev.yml up
```
In this mode, all third-party containers will be deployed. The ports will be forwarded to localhost.
Next run the command
```bash
go test -cover -v ./...
```

### Detailed description of coverage

To get a detailed description of test coverage, run the commands
```bash
go test -v -coverprofile cover.out ./...
go tool cover -html cover.out -o cover.html 
```
Open the `cover.html` file. It will contain fragments of code covered and not covered by tests.
This method works with both unit tests and integration tests.

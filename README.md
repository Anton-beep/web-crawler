# Production

To run a project in production mode

- Set up environment variables in the `configs/Docker.env.template` file
- Install [Docker](https://www.docker.com)
- Run the command
```bash
docker compose -f .\deployments\docker-compose.yml up
```

# Dev

## Pre-setting

To run a project in development mode
- Set up environment variables in the `configs/.env.template` file, copy it to `configs/.env`
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

Чтобы запустить тестирование unit tests, запустите команду
```bash
go test -cover -v ./...
```
В консольле буду выведены результаты тестиования, а также процент покрытия тестов

### Integration tests

Чтобы запустить интеграционные тесты,
- Настройте переменные окружения в файле `configs/.env.template`, скопируйте его в `configs/.env`
- Установите значение переменной `IS_RAN_IN_DOCKER` на `true`.
- Установите [Docker](https://www.docker.com)
- Запустите команду
```bash
docker compose -f .\deployments\docker-compose-dev.yml up
```
В этом режиме будут развёрнуты все сторонние контейнеры. Порты буду прокинуты на localhost.
Далее запустите команду
```bash
go test -cover -v ./...
```

### Подробное описание покрытия

Чтобы получить развёрнутое описание покрытия тестов, запустите команды
```bash
go test -v -coverprofile cover.out ./...
go tool cover -html cover.out -o cover.html 
```
Откройте файл `cover.html`. В нём будут приведены фрагменты кода, покрытого и не покрытого тестами.
Данный метод работает как с unit tests, так и с интеграционными тестами.

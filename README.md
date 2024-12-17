# :globe_with_meridians: Web Crawler

---

<div align="center">

<img src="frontend/public/web_crawler_logo.svg" alt="Web Crawler Logo" width="200" height="200">

### Branch status

| Branch | Pipeline status                                                                                                                                          | Coverage                                                                                                                                                    |
|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `main` | [![Build Status](https://gitlab.crja72.ru/gospec/go1/web-crawler/badges/main/pipeline.svg)](https://gitlab.crja72.ru/gospec/go1/web-crawler/-/pipelines) | [![Coverage Status](https://gitlab.crja72.ru/gospec/go1/web-crawler/badges/main/coverage.svg)](https://gitlab.crja72.ru/gospec/go1/web-crawler/-/pipelines) |
| `dev`  | [![Build Status](https://gitlab.crja72.ru/gospec/go1/web-crawler/badges/dev/pipeline.svg)](https://gitlab.crja72.ru/gospec/go1/web-crawler/-/pipelines)  | [![Coverage Status](https://gitlab.crja72.ru/gospec/go1/web-crawler/badges/dev/coverage.svg)](https://gitlab.crja72.ru/gospec/go1/web-crawler/-/pipelines)  |

Web Crawler is a scalable and efficient tool for collecting information from websites. It supports crawling web pages,
extracting their content, and following links to other pages. Designed with a microservice architecture for performance
and flexibility.


<img src="images/dashboard.png" alt="Dashboard screenshot" width="600" height="auto">
</div>

---

## How it works

```mermaid
flowchart TD
    user(<img src='https://www.svgrepo.com/show/532363/user-alt-1.svg'/> User)
    nginx(<img src='https://www.svgrepo.com/show/303554/nginx-logo.svg'/>)
    frontend(<img src='https://www.svgrepo.com/show/355190/reactjs.svg'/> Frontend static)
    receivers(<img src='https://www.svgrepo.com/show/353795/go.svg'/> Receivers)
    kafka(<img src='https://www.svgrepo.com/show/353951/kafka-icon.svg'/> Kafka)
    collectors(<img src='https://www.svgrepo.com/show/353795/go.svg'/> Collectors)
    analysers(<img src='https://www.svgrepo.com/show/353795/go.svg'/> Analysers)
    redis(<img src='https://www.svgrepo.com/show/354272/redis.svg'/> Redis)
    postgres(<img src='https://www.svgrepo.com/show/354200/postgresql.svg'/> PostgreSQL)
    user -->|request to process a link| nginx
    nginx -->|results of the processing| user
    nginx -->|request to process a link| receivers
    receivers -->|write user data| postgres
    receivers -->|send order to process a link| kafka
    receivers -->|send user data and results of the processing| nginx
    kafka -->|send order to process a link| collectors
    kafka -->|order to analyse the collected text| analysers
    redis -->|collected texts and links| collectors
    postgres -->|texts| analysers
    postgres -->|get user data and graph and analysis data| receivers
    nginx -->|ui| frontend
    frontend -->|ui| nginx
    analysers -->|analysis of a text| postgres
    collectors -->|collected data, links and texts| redis
    collectors -->|formed data, ready for user to read| postgres
    collectors -->|order to collect new links| kafka

```

### Components

**_Nginx_**: Balancer of requests between `frontend` and `receivers` (there can be multiple receivers).

**_Frontend_**: Web interface for users to interact with the system. Notice that the `frontend` here is static files
which are hosted by `Nginx`. `Frontend` is written in React using Vite.

**_Receivers_** (API-Gateway): Service which interacts with user. Receives requests from the `frontend` (actually it's a
http request from the user) and sends them to the `Kafka`. Moreover, handles some user data, i.e. login, registration,
profile update. Authentication is done by using [JWT tokens](https://golang-jwt.github.io/jwt/), users' passwords are
stored in hash.

**_Kafka_**: Service which organizes the work of `analysers`, `receivers`, and `collectors`.

**_Collectors_**: Program which does the actual job in our application. `Collector` crawls around the web and collects
texts and new links.

**_Analyser_**: Program which analyses the collected by a `collector` text and gives some results of the analysis.
`Analyser` can perform various algorithms.

**_PostgeSQL_**: Database.

**_Redis_**: Redis provides storage for `collectors` when they are crawling around the web, so the write and read
operations will be very fast.

### API Documentation

Documentation is now available at `docs/api/apiDocumentation.html` (just open it in your browser).

To rebuild the API documentation, run the following command:

```bash
npx @redocly/cli build-docs api/openapi.json -o docs/api/apiDocumentation.html
```

## Run

### Production

Install [Docker](https://www.docker.com)

#### Using Default Configuration

Use start script (choose yes for default .env variables):

- If you are using Linux:
   ```shell
   bash scripts/start.sh
   ```

- If you are using Windows:
   ```shell
   .\scripts\start.bat
   ```

The web interface will be available at `http://localhost:85`.

#### Using custom .env files

1. Create your .env files in `configs/Docker.env` and `frontend/.env`, you can use `configs/Docker.env.template` and
   `frontend/.env.template` as templates. After that, copy `configs/Docker.env` to `configs/.env`.
2. Run the following command:
   ```bash
   docker compose -f deployments/docker-compose.yml up --build
   ```

### Development

#### Pre-setup

To run the project in development mode:

1. Configure environment variables in `configs/.env` (use `configs/.env.template` as a template).
2. Install [Docker](https://www.docker.com/).
3. Run the following command:
    ```bash
    docker compose -f deployments/docker-compose-dev.yml up --build -V
    ```

In this mode, all third-party containers will be deployed with ports forwarded to `localhost`.

#### Setup

1. Install [Go](https://golang.org/).
2. Install dependencies:
    ```bash
    go mod download
    ```
3. Install [Npm](https://nodejs.org/en/download/package-manager)
4. Create `frontend/.env` file using `frontend/.env.template` as a template. Change `VITE_RECEIVER_HOST` to
   `http://localhost:8080/api`.
5. To run frontend:
    ```bash
    cd frontend
    npm install
    npm run dev
    ```
6. To run backend service run the following command:
    ```bash
    go run cmd/<path_to_service>/<service>.go
    ```

---

### Testing

#### Unit Tests

To run backend unit tests, execute the following command:

```bash
go test ./... -v
```

To run frontend unit tests, execute the following command:

```bash
cd frontend
npm run test
```

#### Integration Tests

1. Start the project in development mode.
2. Run the following command:
    ```bash
    go test ./... -v
    ```

To view detailed test coverage:

```bash
go test ./... -coverprofile coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Open `cover.html` to view the detailed coverage report.

---

### Linting

#### Backend

To lint the backend, install [golangci-lint](https://golangci-lint.run/usage/install/):

```bash
golangci-lint run -c configs/.golangci.yml
```

#### Frontend

To lint the frontend, run the following command:

```bash
cd frontend
npm run lint
```

---

### Documentation

#### API Documentation

To build the API documentation, run the following command:

```bash
npx @redocly/cli build-docs api/openapi.json -o docs/api/apiDocumentation.html
```

API documentation is now available at `docs/api/apiDocumentation.html`

---

### Folder Structure

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

# ThreeDice — Spring Boot

The Java implementation. See the [root README](../README.md) for the game rules, odds table and API reference,
which both services share.

## Tech Stack

- Java 21
- Spring Boot 3.4
- Maven
- Spring Data JPA / Hibernate
- PostgreSQL 16
- Flyway (schema migration)
- Lombok
- Testcontainers + REST Assured (integration testing)

## Prerequisites

- **Java 21**
- **Maven 3.9+**
- **Docker** (for PostgreSQL, and for Testcontainers during tests)

## Running

### 1. Start the database

From the repository root:

```bash
docker compose up -d postgres
```

This starts PostgreSQL 16 on port `5432` with the `threedice` database.

### 2. Run the application

```bash
mvn spring-boot:run
```

The application starts on `http://localhost:8080`. Flyway creates the schema, and a demo player (`id=1`, balance
`1000.00`) is seeded on first startup. The frontend is served at [http://localhost:8080](http://localhost:8080).

### 3. Run the tests

Tests use **Testcontainers** — they spin up their own PostgreSQL container. You only need Docker running.

```bash
mvn clean verify
```

29 integration tests, matched case for case by the [Go suite](../go/README.md).

> The pom pins `testcontainers.version` above the Spring Boot BOM's default. The BOM's 1.20.5 ships a docker-java
> client that cannot negotiate with Docker Engine 29+, which raised `MinAPIVersion` to 1.40; every test fails at
> Docker discovery without the override.

## Configuration

`src/main/resources/application.yml` holds the defaults. Spring Boot's relaxed binding means any of them can be
overridden by environment variable, which is how docker-compose and Coolify configure it:

| Variable | Default |
|---|---|
| `SPRING_DATASOURCE_URL` | `jdbc:postgresql://localhost:5432/threedice` |
| `SPRING_DATASOURCE_USERNAME` | `threedice` |
| `SPRING_DATASOURCE_PASSWORD` | `threedice` |
| `SERVER_PORT` | `8080` |
| `GAME_INITIAL-BALANCE` | `1000.00` |

## Project Structure

```
src/main/java/com/casino/threedice/
├── controller/         # REST controllers (BetController, PlayerController, TransactionController)
├── dto/                # Request/response records
├── entity/             # JPA entities (Client, Bet, Draw, Transaction)
├── exception/          # Custom exceptions & global handler
├── repository/         # Spring Data JPA repositories
├── service/            # Business logic (BetService, PlayerService, TransactionService, DiceEngine)
├── validation/         # Custom Jakarta validation (ValidPredictedValue)
└── ThreeDiceApplication.java

src/main/resources/
├── static/             # Frontend (HTML, CSS, JS)
├── application.yml
└── db/migration/       # Flyway SQL migrations

src/test/java/          # Integration tests (Testcontainers + REST Assured)
```

## Deploying to Coolify

Create an application pointing at this repository:

| Setting | Value |
|---|---|
| Build Pack | `Dockerfile` |
| Base Directory | `/spring-boot` |
| Dockerfile Location | `/Dockerfile` |
| Ports Exposes | `8080` |

Coolify sets the build context to the Base Directory, so the Dockerfile here is self-contained — it only copies
`pom.xml` and `src`.

> **Base Directory must be `/spring-boot`, not `/`.** `Dockerfile Location` is resolved relative to it, so
> `Base Directory: /` + `Dockerfile Location: /spring-boot/Dockerfile` finds the file but builds it with the
> repo root as context, and `COPY src ./src` then fails with `"/src": not found`. Coolify also injects every
> environment variable as a build `ARG`, so expect a `SecretsUsedInArgOrEnv` lint warning for the datasource
> password; it is harmless, and unticking "Build Variable?" on those variables silences it.

Set these environment variables to point at your PostgreSQL service:

```
SPRING_DATASOURCE_URL=jdbc:postgresql://<postgres-host>:5432/threedice
SPRING_DATASOURCE_USERNAME=<user>
SPRING_DATASOURCE_PASSWORD=<password>
```

The Go service deploys as a **second** Coolify application from the same repository with Base Directory `/go`.
It needs its own database on the same PostgreSQL instance — see the [Go README](../go/README.md).

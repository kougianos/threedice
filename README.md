# ThreeDice

#### Input requirements:
<img src="requirements.png" alt="ThreeDice Screenshot" width="500">

#### App screenshot:
<img src="screenshot1.png" alt="ThreeDice Screenshot" width="500">

A web-based dice betting game, implemented **twice** — once in Spring Boot and once in Go — against the same
requirements, the same business logic, the same REST contract and the same frontend.

Players bet on predicting the product of three six-sided dice. If the prediction matches the rolled product, the
player wins their stake multiplied by the applicable odds.

## The two implementations

| | [`spring-boot/`](spring-boot/) | [`go/`](go/) |
|---|---|---|
| Language | Java 21 | Go 1.26 |
| HTTP | Spring MVC | chi + net/http |
| Persistence | Spring Data JPA / Hibernate | sqlc over pgx |
| Money | `BigDecimal` | shopspring/decimal |
| Validation | Jakarta Bean Validation | go-playground/validator |
| Migrations | Flyway | golang-migrate |
| Logging | SLF4J / Logback | log/slog |
| Errors | `ProblemDetail` (RFC 7807) | hand-rolled RFC 7807 |
| Tests | Testcontainers + REST Assured | testcontainers-go + testify |
| Frontend | served from `resources/static` | same files, embedded via `embed.FS` |
| Database | `threedice` | `threedice_go` |
| Port | 8080 | 8081 |

Both own their schema independently, so either can be run on its own against an empty database.

See [`spring-boot/README.md`](spring-boot/README.md) and [`go/README.md`](go/README.md) for building, running and
deploying each one.

## Quick start

```bash
docker compose up --build
```

That starts PostgreSQL 16 (creating both databases), then builds and runs both services:

- Spring Boot → [http://localhost:8080](http://localhost:8080)
- Go → [http://localhost:8081](http://localhost:8081)

Each seeds its own demo player (`id=1`, balance `1000.00`) on first startup and serves the frontend at `/`.

> **Existing volume?** The `threedice_go` database is created by `docker/init-db.sh`, which PostgreSQL only runs
> when the data volume is first initialised. If you already have a `pgdata` volume from an earlier version, either
> `docker compose down -v` (destroys data) or create it by hand:
> ```bash
> docker compose exec postgres psql -U threedice -d threedice -c 'CREATE DATABASE threedice_go OWNER threedice;'
> ```

To run just the database and start a service yourself, `docker compose up -d postgres`.

## Odds Table

| Product Value | Odds |
|---------------|------|
| < 9           | 2x   |
| 9 – 119       | 5x   |
| 120 – 216     | 2x   |

> **Note:** The current odds table heavily favours the house. Since the player must predict the **exact** product of three dice, individual product probabilities are very low (typically 1–7%), yet the maximum payout is only 5x. This results in a large house edge across all ranges - for example, the most common product (e.g. 12 or 24) appears in ~7% of rolls, meaning a fair payout would be around 14x rather than 5x.

### Proposed Fair Odds Table

A more balanced alternative with finer-grained ranges, accounting for actual product probabilities. With 216 total outcomes (6^3), products in the mid-low range (e.g. 12, 24, 30) have the most dice combinations, while extreme values are much rarer.

| Product Value | Proposed Odds | Rationale                                                                 |
|---------------|---------------|---------------------------------------------------------------------------|
| 1 – 8         | 35x           | Only 7 possible products, most with very few combinations (1–9 ways)     |
| 9 – 50        | 12x           | Most populated range - products like 12, 24, 30 each have 12–15 ways    |
| 51 – 120      | 25x           | Fewer combinations per product; values like 60, 72, 120 have 6–9 ways   |
| 121 – 216     | 50x           | Very rare - only 4 products exist (125, 144, 150, 180, 216) with 1–3 ways each |

These multipliers include a ~5–10% house edge over mathematically fair payouts.

## API Reference

Both services expose exactly this contract.

### Place a Bet

```
POST /api/bets
Content-Type: application/json

{
    "playerId": 1,
    "stake": 10.00,
    "predictedValue": 12,
    "idempotencyKey": "b1f3c8e2-..."
}
```

**Response** (201 Created):
```json
{
    "betId": 1,
    "dieOne": 1,
    "dieTwo": 3,
    "dieThree": 4,
    "productValue": 12,
    "status": "WON",
    "winnings": 50.00,
    "balanceAfter": 1040.00
}
```

### Get Player Info

```
GET /api/players/{playerId}
```

**Response** (200 OK):
```json
{
    "id": 1,
    "username": "player1",
    "balance": 1000.00,
    "createdAt": "2026-03-02T09:00:00Z"
}
```

### Bet History (Last 10)

```
GET /api/bets/history/{playerId}
```

**Response** (200 OK):
```json
[
    {
        "betId": 1,
        "predictedValue": 12,
        "stake": 10.00,
        "status": "WON",
        "dieOne": 1,
        "dieTwo": 3,
        "dieThree": 4,
        "productValue": 12,
        "createdAt": "2026-03-02T09:01:00Z"
    }
]
```

### Transaction History (Last 10)

```
GET /api/transactions/history/{playerId}
```

**Response** (200 OK):
```json
[
    {
        "transactionId": 1,
        "betId": 1,
        "type": "CREDIT",
        "amount": 50.00,
        "balanceAfter": 1040.00,
        "createdAt": "2026-03-02T09:01:00Z"
    }
]
```

### Errors

Both services return RFC 7807 `application/problem+json`:

```json
{
    "type": "about:blank",
    "title": "Not Found",
    "status": 404,
    "detail": "Player not found with ID: 99999",
    "instance": "/api/players/99999"
}
```

| Status | When |
|--------|------|
| 400 | Validation failure, insufficient balance, unparseable body or path variable |
| 404 | Unknown player or unmapped route |
| 405 | Known route, wrong verb |
| 409 | Duplicate `idempotencyKey` |

## Features

- **Idempotent Bet Placement**: Each bet request carries a unique idempotency key. Duplicate submissions are detected and rejected with a `409 Conflict`, preventing accidental double-charges.

- **Input Validation**: Stake must be between $1 and $10,000, and the predicted value must be a mathematically possible product of three dice (only 40 out of 216 values are valid). Invalid requests return a `400 Bad Request` with a descriptive error message.

- **Pessimistic Locking on Balance**: Player balance reads use `SELECT ... FOR UPDATE` to prevent concurrent bets from reading stale balances. This ensures correctness under concurrent load.

- **Transactional Integrity**: Each bet placement (stake deduction, dice roll, outcome resolution, balance update, and transaction recording) executes within a single database transaction.

- **Quick Pick & Quick Stake Buttons**: The frontend provides popular predicted values (with win probabilities) and stake increment buttons for faster gameplay without manual input.

- **Secure Dice Rolling**: Dice are rolled from a cryptographically secure source (`java.security.SecureRandom` / `crypto/rand`), ensuring unpredictable outcomes.

- **Transaction Audit Trail**: Every bet generates a corresponding ledger entry (DEBIT on loss, CREDIT on win) with a `balanceAfter` snapshot, exposed via a dedicated history endpoint.

- **Database-Level Constraints**: PostgreSQL CHECK constraints enforce die values (1–6), bet status (WON/LOST), and transaction type (DEBIT/CREDIT) as a defense-in-depth layer beyond application validation.

- **Integration Tests with Deterministic Dice**: Both suites inject a scripted dice source to control outcomes deterministically, running against a real PostgreSQL via Testcontainers. 29 tests each, case for case.


## Measured comparison

Two implementations of the same thing, against the same PostgreSQL, are worth measuring rather than
guessing about. Everything below was measured locally with [`benchmark/`](benchmark/) — rerun it rather
than taking it on trust.

The headline: **Go wins decisively on everything deployment-shaped, and roughly ties on the actual business
operation.**

### Footprint and startup

| | Spring Boot | Go | |
|---|---:|---:|---|
| Docker image | 155.5 MB | 8.9 MB | 17× smaller |
| Deployable artefact | 51.5 MB fat jar *(+ a JRE)* | 14.3 MB static binary *(self-contained)* | |
| **Startup** — restart → first HTTP 200, mean of 6 | **13,983 ms** | **182 ms** | **77× faster** |
| *(spread)* | *12,692 – 14,690 ms* | *167 – 195 ms* | |
| CPU burned booting | 42.82 CPU-s | 0.14 CPU-s | 300× less |
| **Memory, idle** | **313.9 MB** | **5.0 MB** | **63× less** |
| Memory, under load | 389 – 493 MB | 12.9 – 14.7 MB | ~30× less |
| Test suite | 29 tests, ~48 s | 37 tests, ~11 s | |

The ~10 MB the Go service reports in production matches the 5 MB idle / 13 MB loaded measured here. Spring's
314 MB is what an untuned JVM does when given no memory limit — it sizes the heap against available RAM. Under a
container limit it will use less, so that number reflects defaults as much as the JVM itself.

### Warmup: the JVM needs ~50,000 requests

The single most striking difference, and the one that shaped every other reading. Throughput and CPU per request,
measured in consecutive 10,000-request batches from a cold start:

| batch | Spring Boot | | Go | |
|---|---:|---:|---:|---:|
| 1 | 832 req/s | 7.987 ms cpu | 4,480 req/s | 0.525 ms cpu |
| 2 | 1,323 req/s | 4.655 ms | 4,353 req/s | 0.532 ms |
| 3 | 1,742 req/s | 3.122 ms | 4,493 req/s | 0.518 ms |
| 4 | 1,771 req/s | 2.908 ms | 4,053 req/s | 0.545 ms |
| 5 | 2,662 req/s | 1.437 ms | 4,370 req/s | 0.537 ms |
| 6 | 2,710 req/s | 1.359 ms | 4,618 req/s | 0.508 ms |
| … 14 | 2,719 req/s | 1.270 ms | 4,332 req/s | 0.520 ms |

Spring takes roughly **50,000 requests** to reach steady state, and arrives 3.3× faster and 6.3× cheaper than it
started. Go is at full speed on request one and never moves.

This matters beyond benchmarking: a service that scales to zero, redeploys often, or serves bursty traffic may
spend much of its life on the left of that table, never reaching the numbers a steady-state benchmark reports.

It also means **any benchmark that skips warmup is measuring the JIT, not the code**. Judged at batch 1, Spring
looks 5.4× slower and 15× more CPU-hungry than it actually is once warm — and the write path is worse, starting at
16.4 ms of CPU per bet and settling at 3.6 ms. Every figure here was taken after the plateau, which is why the
gaps below are narrower than a naïve run would suggest.

### Read path — `GET /api/players/1`

Steady state, each service measured with **the other one stopped** so they cannot compete for the host's CPUs.
Warmed to plateau, then the median of 5 runs of 5,000 requests at 50 concurrent.

| | Spring Boot | Go | |
|---|---:|---:|---|
| Throughput | 2,746 req/s | 4,069 req/s | 1.5× |
| *(spread over 5 runs)* | *2,127 – 2,964* | *3,833 – 4,734* | |
| Latency p50 | 15.53 ms | 10.45 ms | |
| Latency p95 | 35.74 ms | 21.86 ms | |
| Latency p99 | 60.28 ms | 53.73 ms | |
| **CPU per request** | **1.309 ms** | **0.569 ms** | **2.3× less** |
| RSS during | 389.2 MB | 14.7 MB | |

Throughput on this host is noisy — the spread above is ±15%, and a run taken while the other service was busy
came out 43% higher. CPU per request is the far steadier metric (±8%), because cgroup accounting isolates the
server's own cost from whatever else the machine is doing. Weigh it accordingly.

### Write path — `POST /api/bets`

The real business operation, and the interesting result: **once Spring is warm, it is close to a tie.**

2,000 bets after a 12,000-bet warmup — Spring's write path plateaus at roughly 10,000 bets, Go's is flat from the
first.

At 10 concurrent, all contending for the same player row:

| | Spring Boot | Go | |
|---|---:|---:|---|
| Throughput | 174 req/s | 197 req/s | 1.1× |
| Latency p50 | 59.39 ms | 27.57 ms | Go 2.2× better |
| Latency p99 | 82.30 ms | 585.24 ms | **Spring 7× better** |
| Latency max | 108.75 ms | 2,457 ms | |
| CPU per bet | 3.702 ms | 3.277 ms | 1.1× less |

At 1 concurrent, with no lock contention:

| | Spring Boot | Go | |
|---|---:|---:|---|
| Throughput | 139 req/s | 144 req/s | 1.04× |
| Latency p50 | 6.87 ms | 6.48 ms | |
| **CPU per bet** | **2.494 ms** | **3.014 ms** | **Spring 1.2× less** |

Three things worth drawing out, none of which flatter the obvious narrative:

- **The write path is database-bound, not language-bound.** Each bet is five statements in one transaction behind
  `SELECT ... FOR UPDATE`. PostgreSQL sets the pace, so Go's advantage collapses to ~1.1×. Rewriting in Go does
  not make the row lock any faster.
- **Uncontended, the warmed JVM uses *less* CPU per bet than Go** (2.494 ms vs 3.014 ms). C2 is genuinely good at
  this, and the honest reading is that Go's runtime is not what wins the deployment-shaped comparison above —
  its lack of a warmup phase and its memory floor are.
- **The two fail differently under contention.** Go has the far better median (27.57 ms vs 59.39 ms) but a much
  worse tail (p99 585 ms vs 82 ms). That is a configuration artefact, not a language one: pgx defaults its pool to
  `max(4, NumCPU)` = 8 while HikariCP defaults to 10, so against 10 concurrent clients Spring's pool fits exactly
  and Go's leaves two requests queueing every round. Neither pool is tuned. This is what defaults do.

### Reproducing

```bash
docker compose up -d --build
cd benchmark

# Stop the other service first -- otherwise they compete for the same CPUs.
docker compose stop go

# Warm properly. Five batches minimum for the JVM, or you are timing the JIT.
for i in 1 2 3 4 5; do go run . -url http://localhost:8080 -mode get -n 10000 -c 50 -quiet; done
go run . -url http://localhost:8080 -mode get -n 5000 -c 50

docker compose start go && docker compose stop spring-boot
go run . -url http://localhost:8081 -mode get -n 5000 -c 50   # Go needs no warmup
```

CPU and memory come from each container's own cgroup accounting, which is exact, rather than from `docker stats`
sampling — which reported 0% for a container that was demonstrably busy:

```bash
docker compose exec go sh -c 'cat /sys/fs/cgroup/cpu.stat; cat /sys/fs/cgroup/memory.current'
```

Take a reading before and after a fixed number of requests; the delta is the CPU actually spent on that work.

### What these numbers are not

- Measured on one machine — Docker Desktop on Windows, 8 CPUs, 4 GB — with **the load generator sharing those same
  CPUs** as the service under test. That caps the absolute throughput figures for both. Treat the ratios as
  indicative and the absolutes as a floor; CPU-per-request is the more robust metric, since cgroup accounting
  isolates the server's own cost.
- **Neither service is tuned.** No JVM flags, no heap sizing, no pool sizing, no GC selection. A tuned JVM with a
  sensible heap cap and a tuned pgx pool would both look different.
- This measures *this* application: a thin CRUD service over PostgreSQL that spends most of its wall time waiting
  on a row lock. It is not a general Java-versus-Go benchmark, and the write-path result is the honest
  illustration of why.
## Known differences between the two services

The APIs were diffed response-by-response; 22 of 23 error and validation cases are byte-identical. The
remainder are deliberate:

- **Money always renders at 2dp in Go.** `BigDecimal` carries the scale of whatever it was built from, so Java
  answers a `"stake": 10` request with `"winnings": 50` but a `"stake": 10.00` request with `"winnings": 50.00`,
  and a loss with `"winnings": 0`. Go always emits `50.00` / `0.00`. Every value read from the database is 2dp in
  both. The same applies inside the insufficient-balance message.
- **`createdAt` within a single bet.** Java stamps each entity with `Instant.now()`, so a bet and its transaction
  differ by a few hundred microseconds. Go lets PostgreSQL default the column, and `now()` is the transaction
  timestamp, so both rows share one value. Ordering across bets is unaffected.

## Further improvements that could be delivered in future iterations

- Bet and transaction history to return paginated results, not only last 10 entries.
- Add filter capabilities in bet and transaction history.
- Add cache layer for hot data.
- Add sign up / login flows and authenticate bets with JWT through spring security.
- Add some animations and assets (eg rolling dice) on the client side and improve UX.
- The frontend is duplicated in both services; a shared build step would stop the two copies drifting apart.

# ThreeDice

#### Input requirements:
<img src="requirements.png" alt="ThreeDice Screenshot" width="500">

#### App screenshot:
<img src="screenshot1.png" alt="ThreeDice Screenshot" width="500">

A web-based dice betting game built with Spring Boot 3, PostgreSQL, and Java 21.

Players bet on predicting the product of three six-sided dice. If the prediction matches the rolled product, the player wins their stake multiplied by the applicable odds.

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

## Prerequisites

- **Java 21**
- **Docker** (for PostgreSQL or for Testcontainers during tests)
- **Maven 3.9+**

## Getting Started

### 1. Start the database

```bash
docker compose up -d
```

This starts a PostgreSQL 16 instance on port `5432` with database `threedice`.

### 2. Run the application (maven and Java 21 required on PATH)

```bash
mvn spring-boot:run
```

The application starts on `http://localhost:8080`. Flyway automatically creates the schema, and a demo player (`id=1`, balance `1000.00`) is seeded on first startup.

The **frontend** is served automatically at [http://localhost:8080](http://localhost:8080)

### 3. Run the tests

Tests use **Testcontainers** - they spin up their own PostgreSQL container automatically. You only need Docker running.

```bash
mvn clean test
```

## API Reference

### Place a Bet

```
POST /api/bets
Content-Type: application/json

{
    "playerId": 1,
    "stake": 10.00,
    "predictedValue": 12
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

## Features

- **Idempotent Bet Placement**: Each bet request carries a unique idempotency key. Duplicate submissions are detected and rejected with a `409 Conflict`, preventing accidental double-charges.

- **Input Validation**: Requests are validated using Jakarta Bean Validation annotations. Stake must be between $1 and $10,000, and the predicted value must be a mathematically possible product of three dice (only 40 out of 216 values are valid). Invalid requests return a `400 Bad Request` with a descriptive error message.

- **Pessimistic Locking on Balance**: Player balance reads use `SELECT ... FOR UPDATE` to prevent concurrent bets from reading stale balances. This ensures correctness under concurrent load.

- **Transactional Integrity**: Each bet placement (stake deduction, dice roll, outcome resolution, balance update, and transaction recording) executes within a single database transaction.

- **Quick Pick & Quick Stake Buttons**: The frontend provides popular predicted values (with win probabilities) and stake increment buttons for faster gameplay without manual input.

- **Secure Dice Rolling**: Dice are rolled using `java.security.SecureRandom`, ensuring cryptographically unpredictable outcomes .

- **Transaction Audit Trail**: Every bet generates a corresponding ledger entry (DEBIT on loss, CREDIT on win) with a `balanceAfter` snapshot, exposed via a dedicated history endpoint.

- **Database-Level Constraints**: PostgreSQL CHECK constraints enforce die values (1–6), bet status (WON/LOST), and transaction type (DEBIT/CREDIT) as a defense-in-depth layer beyond application validation.

- **Integration Tests with Deterministic Dice**: Tests use `@MockitoBean` on `DiceEngine` to control outcomes deterministically, running against a real PostgreSQL via Testcontainers.

## Tech Stack

- Java 21
- Spring Boot 3.4
- Maven
- Spring Data JPA / Hibernate
- PostgreSQL 16
- Flyway (schema migration)
- Lombok
- Testcontainers + REST Assured (integration testing)

## Further improvements that could be delivered in future iterations

- Bet and transaction history to return paginated results, not only last 10 entries.
- Add filter capabilities in bet and transaction history.
- Add cache layer for hot data.
- Add sign up / login flows and authenticate bets with JWT through spring security.
- Add some animations and assets (eg rolling dice) on the client side and improve UX.

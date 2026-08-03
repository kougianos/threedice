# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Stack

Vanilla HTML, CSS and JavaScript, no framework, no bundler, no Node toolchain. The pages are static files served
straight from each service (`go/web/static/` embedded via `embed.FS`, `spring-boot/src/main/resources/static/`).

Delegated: the user left the delivery decision to design work. Chosen constraint is to keep the frontend
dependency-free, because the repository's whole argument is a service you can build and run with only Docker plus
one language toolchain, and a reviewer who has to install Node before seeing the game pays a cost the UI cannot
earn back. Chosen direction is one canonical copy of the frontend, copied into both services at build time, so the
two cannot drift. **That de-duplication is not implemented yet**: until it lands, every frontend edit must be
applied to both copies, which are currently byte-identical.

## Users

The primary visitor is a hiring reviewer or interviewer assessing the author's engineering craft. They arrive from
the repository, run `docker compose up --build`, open `localhost:8080` or `localhost:8081`, place a handful of bets
to confirm the thing works end to end, and then go back to reading code, tests and the README. They are evaluating,
not gambling. Their session is short, deliberate, and they know exactly what the requirements document asked for.

No other audience is confirmed. There are no real players and no real money.

## Product Purpose

ThreeDice is a dice betting game built to a take-home specification (`requirements.png`). A player stakes an amount
on the product of three six-sided dice; if the roll's product matches the prediction, they win stake times the
applicable odds, and their balance updates.

The specification called the frontend optional. It exists anyway, so the submission can be *used* rather than only
read, and success is that a reviewer forms a favourable judgement of the engineering behind it within a minute of
first paint.

## Positioning

The same requirements implemented twice, in Spring Boot and in Go, against one REST contract and one frontend, then
measured rather than argued about. A neighbouring take-home submission cannot truthfully claim two production-shaped
implementations in parity, 22 of 23 error responses byte-identical between them, and a benchmark harness in the repo
that reproduces every figure quoted.

## Operating Context

- `docker compose up --build` starts PostgreSQL 16 and both services. Spring Boot on 8080, Go on 8081, each serving
  the frontend at `/` and seeding its own demo player on first startup.
- Each service owns its schema independently (Flyway / golang-migrate) against its own database (`threedice`,
  `threedice_go`), so either runs alone.
- The root `README.md` is the primary artefact for this audience and carries the game rules, the API reference, the
  known differences between the two services, and the measured comparison. The UI is read alongside it.
- `benchmark/` is a Go load generator used for the published numbers.
- Both suites use Testcontainers against real PostgreSQL with a scripted dice source for determinism.
- Deployment beyond `docker compose` is unsettled. Earlier Coolify configuration was committed and reverted.

## Capabilities and Constraints

**Frozen.** The REST contract documented in the README is fixed, and both services must keep answering it
identically. This is the one commitment the user declared explicitly.

**Current facts, not commitments.** The following are true today and were deliberately *not* claimed as immutable.
Treat changing any of them as a product decision to raise with the user, never a silent one:

- Odds are 2x below 9, 5x for 9 to 119, 2x for 120 to 216, as the specification defines them. The README proposes a
  fairer alternative table with its reasoning; that proposal is commentary and is not implemented.
- One seeded demo player (`id=1`, username `player1`, balance `1000.00`), hardcoded as `PLAYER_ID` in the frontend.
  No signup, no login, no authentication, play money only.
- The frontend is duplicated byte-for-byte in both services. See `## Stack` for the chosen direction.

**Rules the UI must respect:**

- Stake is between `1.00` and `10000.00` inclusive.
- Only 40 of the 216 outcomes are reachable as a product of three dice, so a prediction can be `1..216` numerically
  and still be rejected. The current input labels its range as 1 to 216, which is the numeric bound rather than the
  valid set.
- Every bet carries a client-generated `idempotencyKey`; a repeat returns `409 Conflict`.
- Errors arrive as RFC 7807 `application/problem+json`. Validation failures return a joined `detail` string that is
  identical across both services and is the only error text the UI has to work with.
- Bet and transaction history return the last 10 entries only. No pagination, no filtering.
- Money renders at 2dp. Spring may emit `50` where Go emits `50.00`; the frontend must not trust the server's
  formatting.
- Dice are drawn from a cryptographically secure source in both services.

## Brand Commitments

Name: **ThreeDice**. No visual, voice or identity constraints were declared binding. Existing marks in the
repository (`favicon.svg`, present identically in both services) are incumbent artefacts, not commitments.

## Evidence on Hand

Real, and usable:

- `requirements.png` - the actual specification the work was judged against, including the odds table and the
  optional-endpoint list.
- `screenshot1.png` - the incumbent interface as shipped.
- `benchmark/` plus the measured comparison in `README.md` - image size, startup, CPU per request, memory, warmup
  curve, read and write path latency, all reproducible and honestly caveated in the README's own words about what
  the numbers are not.
- Two integration suites running against real PostgreSQL.

Absent, and not to be fabricated: users, testimonials, customers, traffic, revenue, pricing, licensing, uptime,
awards, and any claim of production deployment. There is no real money and no regulated gambling operation behind
this. Do not invent responsible-gambling apparatus, operator branding, or compliance copy to fill the space.

## Product Principles

1. **The reviewer is the user.** Every surface decision answers "does this raise or lower their estimate of the
   engineering?" Decoration that a careful reader would price as a distraction costs more than it returns.
2. **Never claim more than the repository can prove.** The evidence above is unusually strong; padding it with
   invented social proof would poison the part that is real.
3. **Parity is the product.** Anything shipped to one service ships to the other, unchanged, or the central claim
   stops being true.
4. **Honesty about the game itself.** The implemented odds favour the house heavily and the README says so plainly.
   The interface must not paper over that, and must not moralise about it either.
5. **Runnable with what a reviewer already has.** Docker plus one language toolchain. Anything that adds a setup
   step has to justify itself against a visitor who may simply not bother.

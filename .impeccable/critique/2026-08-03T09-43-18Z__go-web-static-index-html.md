---
target: the game surface (go/web/static/index.html)
total_score: 18
max_score: 40
na_heuristics: 
p0_count: 1
p1_count: 3
timestamp: 2026-08-03T09-43-18Z
slug: go-web-static-index-html
---
Method: dual-agent (A: a148aa09374aae748 · B: a3784f11f8dc2d95f)

Target: the game surface, `go/web/static/index.html` plus its stylesheet and script. Reviewed live against a local mock of the REST API on port 8099 (now stopped), because neither service was running and no containers were up. Mode: **Operate**, with PRODUCT.md's hiring-reviewer audience raising the weight of first impression.

## Design Health Score

| # | Heuristic | Score | Key Issue |
|---|-----------|-------|-----------|
| 1 | Visibility of System Status | 2 | The verdict renders ~200ms before the dice that justify it, and the previous bet's result stays on screen for the whole round trip with no pending state |
| 2 | Match System / Real World | 2 | `predictedValue: must be a possible product of three dice` leaks a JSON key at the user; "Predicted Product (1-216)" states a numeric bound as if it were the legal set |
| 3 | User Control and Freedom | 2 | The toast cannot be dismissed or re-read (3s then gone); no repeat-last-bet; a bet is correctly irreversible, which is honest |
| 4 | Consistency and Standards | 2 | Visual system is genuinely consistent; interaction is not. Three error channels, five different focus colors across 16 stops, and `#` means two different things in the two tabs |
| 5 | Error Prevention | 1 | The client validates 1..216 while knowing only 40 values are legal, so it knowingly submits 176 values that always fail. No `min`/`max`/`step`, and over-balance stakes are sent despite the balance being in the header |
| 6 | Recognition Rather Than Recall | 2 | Quick picks are a real strength, but the odds legend sits *after* the submit button, and the result card never shows what you predicted |
| 7 | Flexibility and Efficiency | 2 | Enter submits, quick stakes accumulate rather than overwrite. No repeat-bet, 13 tab stops to reach the primary action |
| 8 | Aesthetic and Minimalist Design | 3 | The strongest layer by far. Undercut by a non-interactive solid-Iris `5x` chip reading as a second button 40px under the real one |
| 9 | Error Recovery | 1 | Errors appear ~400px from their cause, unlabelled, for 3 seconds, with no field state, no focus move, no announcement, and silent loss entirely if the tab is backgrounded |
| 10 | Help and Documentation | 1 | No rules, no mention of the 40-value constraint, no house edge. Per-pick help lives in `title`, unreachable by keyboard or touch. The footer links nowhere |
| **Total** | | **18/40** | **Poor (45%)** |

That band is harsh and I am going to stand behind it, with one qualification: the score is dragged down almost entirely by behavior, not by looks. Heuristic 8 is the only one above 2. If you fixed nothing but the five priority issues below, this lands in the high twenties without touching the visual system.

## Design Specificity Verdict

**LLM assessment: roughly 85% category-interchangeable. One authored idea, executed once, revealed late, and never developed.**

Swap the copy and the shell is a CI dashboard or a crypto portfolio: sticky header with gradient wordmark left and identity pill right, a 900px centered column of shadowed cards, a tabbed "recent activity" module with uppercase tracked headers and tinted status badges, a label-over-input form with quick-fill chips and a full-width violet submit. Violet primary plus teal success plus coral danger on near-black is the most common dark SaaS scheme in existence. DESIGN.md's names for them (Electric Iris, Winner's Teal, Ember Clay) are authorship applied *after* the fact rather than authorship in the choice.

What is genuinely authored is the **dice equation**: three 64px tiles reading `3 × 3 × 5 = 45`, with Fog operators and a gold product. It could only exist for this product, it teaches the game's rule without a word of instruction, and the shared 12px radius between a die and a card is a real system decision that pays off. The **quick picks carrying win probabilities** are the second authored move.

The problem is that the authored idea is invisible when it matters most. The result card is `display: none` until the first bet, so **on a cold load there are no dice on the page at all**. DESIGN.md calls gold-on-slate "the whole system," and the only gold at first paint is a single quick-pick border. A reviewer's first five seconds contain a form and an empty seven-column table.

The largest missed opportunity: **the 40-of-216 fact is the most interesting thing about this game and appears nowhere in the interface.** A picker built from the 40 reachable products, banded by payout, would be specific to this product alone, would teach the odds without a legend, and would make the impossible-prediction error structurally unreachable. Related: the result card never shows what you predicted, so the product's single question ("did these match?") is rendered as an equation with one side missing.

**Deterministic scan.** `detect.mjs --json go/web/static` exited 2 with **14 findings**, all in `go/web/static/css/style.css`: `design-system-font-size` x12 (lines 72, 129, 191, 202, 263, 308, 326, 347, 383, 412, 490, 519), `gradient-text` x1 (line 59), `design-system-color` x1 (line 492).

**12 of those 14 are false positives, and they are my fault, not the code's.** The detector reconciles against DESIGN.md's YAML frontmatter, whose `typography:` block declares only five steps. The prose body documents roughly ten (0.65, 0.8, 0.85, 0.9, 0.95, 1, 1.1, 1.3, 1.5, 1.8, 2rem), so the detector fires on sizes the design system explicitly specifies. Line 59's gradient text and line 492's toast shadow are both documented verbatim in DESIGN.md and are sanctioned. Only two findings are genuine and both are trivial: `.result-status` at 1.3rem (line 347) and the `.die` numeral dropping to 1.4rem under 600px (line 519) are sizes DESIGN.md never records. The scan's real output is a note about the spec I wrote, not about the stylesheet.

**Visual overlays.** Injection succeeded: mutation was available with no CSP obstruction, the live server ran on port 8400, `detect.js` loaded and executed in the page, and the console reported **18 anti-patterns** (19 distinct emitted lines, so its own header count is off by one). Both servers are now stopped and both ports are confirmed closed. **No user-visible overlay is available for you to look at**: the Browser pane was not compositing for the whole session, so screenshots timed out and all visual evidence is computed-style and geometry measurement rather than pixel inspection.

The overlay flagged `undersized-ui-text` x8 (every `.pick-btn small` at 10.4px against an 11px floor), `ai-color-palette` x6 ("cyan neon on dark" on the teal elements), `gpt-thin-border-wide-shadow` x3 (all three cards: 1px border plus 24px blur), `gradient-text` x1, and `low-contrast` x1 on `.pick-btn.best small` at 3.2:1, which reproduced exactly in independent measurement (3.17:1).

Two of those deserve a straight answer. The teal "neon" flag is a false positive: teal is a documented semantic outcome color, not palette drift. The thin-border-wide-shadow flag is **not** a false positive. It is a correct identification of a saturated machine-generated pattern, and DESIGN.md canonizes it as "card ambient." For an audience who has seen a thousand of these cards, that is worth knowing even though the project adopted it deliberately.

**Where the two assessments met.** The detector caught contrast failures the design review missed: `.pick-btn.best small` at 3.17:1, `.badge-lost` and `.badge-debit` at 4.35:1, and `.tab.active` Iris on card at 3.46:1, all below the 4.5:1 they need. The design review caught everything the detector structurally cannot: wrong published probabilities, the verdict/dice race, the stale result card, the direction-blind balance pop, and the reindexing `#` column. Neither would have been sufficient alone.

## Overall Impression

The visual layer is better than this project's audience will give it credit for, and the behavior layer is worse than the backend it sits on. The restraint is real: three static files, zero font requests, zero image requests, one 5-pip SVG favicon, one ambient shadow, no confetti anywhere near a betting UI. That is the right instinct for a reviewer judging engineering, and it is consistently executed.

Then the interface publishes four wrong numbers, tells you that you lost before it shows you the dice, and prints a raw JSON field name at you in a bar that disappears in three seconds. The Go service deliberately stacks its valid-product rule outside the validator tag chain so it can report multiple failures exactly like Bean Validation does. That is careful, and the frontend throws the result away as a string.

**The single biggest opportunity:** the frontend is currently undercutting the backend's best work in front of the exact audience the backend was written to impress.

## What's Working

1. **The dice equation composition.** Rendering the outcome as literal arithmetic makes the rule self-evident: you see *why* the product is what it is with no instruction. The 12px radius shared between dice and cards is a genuine system decision, and it lands: the tiles read as belonging to the layout rather than dropped into it.

2. **Quick stakes add rather than set.** `+$10 / +$100 / +$250` accumulate onto the current value, so two taps build $260. Small, correct, non-obvious, and it matches how people actually think about raising a stake. The teal treatment signals "money in" without stealing the primary action's violet.

3. **Restraint under category pressure.** The pull toward neon, coin-shine and confetti in a betting UI is enormous and this holds a quiet line. Zero console errors, zero failed requests, clean network behavior across a full bet cycle. For the stated audience this is the right register.

## Priority Issues

### [P0] Four of eight quick picks publish wrong win probabilities, and one contradicts the app's own odds table

Computed over all 216 ordered triples and verified independently of the review agent:

| Pick | UI shows | Actual | Combinations |
|---|---|---|---|
| 24 | 5.6% | **6.94%** | 15 |
| 30 | 2.8% | **5.56%** | 12 |
| 36 | 5.1% | **5.56%** | 12 |
| 60 | 3.7% | **5.56%** | 12 |

30 is understated by a factor of two. Separately, `index.html:40` gives pick **6** the tooltip "4.2% chance - 5x payout", but 6 is under 9, so the legend rendered 40px below it pays **2x**. And the star on 12 is not an exclusive truth: 24 has identical combinations (15), identical payout (5x) and identical expected value (0.3472). They are exactly tied, while the interface shows 24 at a falsely lower probability, which makes the false claim look substantiated. Three picks (30, 36, 60) share one true probability and are shown as three different numbers.

**Why it matters:** PRODUCT.md Principle 2 is "Never claim more than the repository can prove," the repo's pitch is "measured rather than argued about," and DESIGN.md's own Do says "pair every quick pick with its real win probability." The one place the UI quotes a measurement, it is wrong, and a reviewer can falsify it in thirty seconds with a calculator. This is a correctness failure, not a taste disagreement, which is exactly the category that costs you the judgement PRODUCT.md defines as success.

**Fix:** correct to `6 → 4.2% / 2x`, `12 → 6.9% / 5x`, `24 → 6.9% / 5x`, `36 → 5.6% / 5x`, `8 → 3.2% / 2x`, `30 → 5.6% / 5x`, `60 → 5.6% / 5x`, `120 → 2.8% / 2x`. Move the star to the tied pair or label both. Better: derive the whole table at load from a five-line triple loop so it can never drift again, which also demonstrates the 40-of-216 fact in code. Mirror into `spring-boot/src/main/resources/static/`.

**Suggested command:** `/impeccable clarify`

### [P1] The verdict renders ~200ms before the dice, and the previous bet's result stays on screen during the roll

`app.js:109-131` writes `statusEl.textContent` synchronously while deferring the die faces by `setTimeout(..., 200)`. Measured frame by frame on the live page:

```
t=  3ms   dice "?,?,?"   status "LOST"   balance $970
t= 49ms   dice "?,?,?"   status "LOST"   balance $960   <- money already gone
t=195ms   dice "?,?,?"   status "LOST"
t=243ms   dice "6,2,6"   status "LOST"   <- faces finally land
```

The 400ms flip, which DESIGN.md calls the signature moment, plays out entirely after the answer has been given. Worse, `showResult` never clears prior state, so from click until the response returns the card displays the *previous* bet's dice, product and verdict with no pending indication. After a win you can sit looking at a WON pill attached to a bet still in flight.

**Why it matters:** this is the only designed moment in the product, and the sequencing destroys it. On a slow link it is not just anticlimactic, it is misleading.

**Fix:** move the `won` / `statusEl` / `winningsEl` block inside the existing `setTimeout` at `app.js:120`. At the top of the submit handler, reset dice and product to `?`, clear `statusEl.className` to a neutral `result-status`, and blank the winnings. Four lines, and it is the highest-value four lines in this report.

**Suggested command:** `/impeccable animate`

### [P1] Three error channels, two registers, and a three-second bar 400px from the cause

Observed from a single form:

| Input | Channel | Message |
|---|---|---|
| stake blank | native browser bubble | "Please fill in this field." |
| stake 0 or "abc" | toast | "Enter a valid stake" |
| stake 0.50 | toast | "stake: must be greater than or equal to 1.00" |
| stake 99999 | toast | "stake: must be less than or equal to 10000.00" |
| stake 5000 on $900 | toast, after a round trip | "Insufficient balance. Required: 5000.00, Available: 900.00" |
| prediction 0 or 217 | toast | "Predicted value must be between 1 and 216" |
| prediction 13 | toast | "predictedValue: must be a possible product of three dice" |

`app.js:71` accepts any integer 1..216 while only 40 are legal, so the UI knowingly submits 176 values that always 400. Nothing on the page explains the constraint. The failure arrives as a raw API string with the camelCase key intact, in a toast that cannot be dismissed, cannot be re-read, self-destructs in 3s, and leaves the offending value styled as perfectly valid. The only error that gets field-adjacent presentation is the native bubble the design system did not author.

And it can vanish entirely: `showToast` (`app.js:254`) adds `.show` inside a single `requestAnimationFrame` while the 3-second removal uses `setTimeout`. With the tab backgrounded, rAF does not fire, so every toast sat 19px below the viewport for its whole life and removed itself unseen. The app's single error channel silently discards errors, with no fallback, because DESIGN.md made the toast the only channel.

**Fix:** generate the 40 legal products client-side and validate before submit ("12 is a possible product; 13 isn't. Try 12, 24 or 36."). Add a real field error state (`border-color: var(--danger)`, `aria-invalid="true"`, `aria-describedby` to an inline message). Pre-check `stake > balance` against the balance already in the header. Strip the `field: ` prefix off server detail strings. Give the toast `role="alert"`, a dismiss control, and persistence for validation errors.

**Suggested command:** `/impeccable harden`

### [P1] The history table overflows the viewport by 63px on mobile, and misinforms on desktop

There are **zero `overflow` declarations in the entire stylesheet**. At 375px with history populated, `document.documentElement.scrollWidth` is **438** against a `clientWidth` of **375**. Hiding the history card drops it to exactly 375, so that card owns all 63px: a 7-column table with a natural width of 395.7px inside a 293px content box, with `table-layout: auto`, `min-width: 0`, and `overflow-x: visible` on every ancestor up to `body`. Result and Time render past the edge of the screen. Rows run 75px tall with cells wrapped to four lines.

On desktop, `Time` is the widest column (174px of 818) for a value whose second-granularity still fails to disambiguate adjacent rows, `#` is leftmost, and `Predicted` is separated from `Product` by `Dice`, so the two numbers whose comparison is the entire point of the row cannot be scanned side by side. In normal play four of seven columns are constant.

The `#` column is a page index, not an identity: it is `bets.length - i` (`app.js:148`), so placing one bet **relabels every row in the column**, while `betId` arrives as the first field of every history entry and is discarded. The same `#` header means bets in one tab and transactions in the other, and the two desynchronise permanently on the first win.

**Fix:** wrap the table in `overflow-x: auto` with a `min-width`, use `betId` / `transactionId`, reorder to `Predicted | Product | Dice | Stake | Result | Time`, add `font-variant-numeric: tabular-nums` with right-aligned numerics, and below 600px collapse each row into a stacked two-line card. While there: the sticky header is 118px tall on mobile, roughly 18% of a phone's usable height for a wordmark and a balance.

**Suggested command:** `/impeccable adapt`

### [P2] Silent to assistive tech, unfocusable-looking for 14 of 16 controls, and three contrast failures

- **Zero live regions on the page.** No `aria-live`, no `role="status"`, no `role="alert"` anywhere. The roll result, the balance change and every toast are announced to nobody. The toast ships as a bare `<div class="toast">` whose only attribute is `class`.
- **The quick picks' accessible names omit the number they select.** `title` wins the name computation, so the button announces as "4.2% chance - 5x payout" with no mention of 6. Two buttons announce near-identically. `title` is also unreachable by keyboard and touch, so it is the wrong home for the only per-pick explanation.
- **The only `:focus` rule in the stylesheet is `.form-group input:focus` (lines 145-149).** All 14 buttons fall back to Chrome's UA ring, giving five different focus colors across one form, none of them the Iris the system says focus means. The one authored ring composites to 1.23:1 against the card and contributes nothing; the real indicator is the border shift to Iris at 3.88:1.
- **Contrast failures:** `.pick-btn.best small` 3.17:1, `.badge-lost` and `.badge-debit` 4.35:1, `.tab.active` 3.46:1, all against a 4.5:1 requirement. Every UI boundary (`#2a2e3e` borders) sits at 1.25 to 1.40:1 against 3:1.
- **No `prefers-reduced-motion` rule exists.** The only media query in the file is `max-width: 600px`, while the signature interaction is a 360 degree Y rotation with a scale dip on every bet.
- No `scope` on any of the 12 `<th>`, no captions, no heading on the history card so heading navigation skips it, and the tabs have no `role="tab"` / `aria-selected` / `aria-controls` / arrow keys.

**Fix:** a visually-hidden `role="status" aria-live="polite"` updated with "Rolled 3, 3, 5. Product 45. Lost. Balance $900."; `role="alert"` on the toast; replace `title` with `aria-label="Predict 6, 4.2% chance, pays 2x"` and surface payout as visible text; one `:focus-visible { outline: 2px solid var(--primary); outline-offset: 2px }` for everything; lift the three failing tints; add the reduced-motion block.

**Suggested command:** `/impeccable audit`

## Persona Red Flags

**Morgan, the hiring reviewer (project-specific, from PRODUCT.md).** Arrives from the repo having read the parity and benchmark claims. Budget: ninety seconds. The page loads with **no dice on it**, so her first five seconds contain a violet form above an empty seven-column table. She reads the quick picks, the only thing teaching the game, knows there are 216 triples, counts 12 combinations for 30 against the 2.8% shown, and finds four of eight wrong. Then she notices pick 6 promises 5x while the legend two inches below says products under 9 pay 2x. She places a bet and the card says LOST before the dice have faces; she notices, because noticing races is the job. She types 13, the exact probe any reviewer runs, and the frontend prints the API's raw string at her. Net: the restraint and the dependency-free delivery raise her estimate; the wrong numbers, the spoiled reveal and the leaked JSON key lower it further, because they are correctness failures.

**Jordan (confused first-timer).** The label says any number 1-216 is legal, so he types 13. No inline warning, no border change. He presses Roll and an orange bar appears 400px below where he was looking reading `predictedValue: must be a possible product of three dice`, and vanishes in 3 seconds. He does not know what `predictedValue` is, and nothing on the page ever explains the 40-value rule, so his only discovery path is trial and error. If he wants a typed value's payout, the legend that would tell him is *below* the button he already pressed. After the roll the result never shows what he predicted, so he cannot check whether he was close.

**Sam (screen reader, keyboard only, needs 4.5:1).** Tabs to the picks and hears "button, 4.2% chance - 5x payout". **No number.** Eight buttons that describe odds and never say what they select. Presses Enter on Roll and **nothing is announced**: dice change, balance changes, history rewrites, all silent. If it errored, still nothing. Reaches the history and finds two unlabelled buttons with no tab semantics, no `scope` on any header, and no heading on the card. Focus is technically visible on all 16 stops but is a different browser-default color on nearly every control. Every bet forces a 360 degree 3D rotation with no reduced-motion opt-out.

**Casey (distracted, one-handed, mobile).** The sticky header eats 118px, roughly 18% of her usable height, for a wordmark and a balance. The quick-stake chips are **26.59px tall**, 6.4px apart, and because they *add* rather than set, a mis-tap silently inflates her stake instead of erroring. She taps Roll and nothing appears to happen: there is no scroll management and the result lands 230px below the button. The history below hangs 63px past the right edge of the screen with no scroll container, in 75px rows of wrapped text.

## Minor Observations

- The bet card carries **two solid-Iris elements**: the submit button and the non-interactive `<span class="odds-badge highlight">5x</span>`. That breaks both the One Loud Thing Rule and the Iris Means Action Rule, and visually the chip reads as a second button 40px under the real one.
- **Lucky Gold has four homes in the shipped CSS** (wordmark gradient, best pick, die numerals, product), not the three its Named Rule declares. DESIGN.md contradicts itself here: the Secondary prose grants the wordmark, the Named Rule does not.
- `.die { transition: transform 0.3s }` (`style.css:311`) is both dead code (the roll uses `animation`) and a breach of the Two-Tenths Rule.
- The balance pop is direction-blind: `updateBalance` always adds `.updated`, so the pill performs the same congratulatory 1.1 scale going 1000 to 1050 and 1000 to 990.
- Money is left-aligned with proportional figures, so `$1040.00` and `$50.00` do not align on the decimal.
- `app.js:4` says "Communicates with the Spring Boot REST API" and `app.js:8` comments `PLAYER_ID` as seeded by `DataInitializer`, a Spring class name, in the file the Go service ships byte-identically. Principle 3 is "Parity is the product."
- The footer reads "ThreeDice - A Coding Challenge" and links nowhere, not even to the README that PRODUCT.md calls the primary artefact for this audience. Cheapest available win for the reviewer.
- The empty state copy is good, but it renders under seven uppercase column headers for a table with no data. Progressive disclosure is inverted.
- Every fix here has to be applied twice by hand, because the frontend is duplicated byte-for-byte.

## Questions to Consider

1. If only 40 of 216 values can ever win, why is the prediction a text box at all? What does this look like if the input *is* the distribution, so "13" is not a thing you can type, the odds are legible without a legend, and the most interesting fact about the game becomes the primary composition?
2. The result card never shows what you predicted. If the product's single question is "did these match," why is the answer rendered as an equation with one side missing?
3. The dice equation is the only authored idea here, and it is invisible until money is committed. What if the cold page opened on a resting `? × ? × ? = ?`, so the first five seconds contained the product's actual signature?
4. Losing is the 93% case and every loss is currently identical. Is there an honest, non-manipulative way to make a loss informative, showing what would have won, without simulating a slot machine's near-miss?
5. DESIGN.md deliberately gives fields no error state and makes the toast the single channel for everything. That one decision is the root of four separate findings here. What was it buying?
6. Two tabs show the same events with no join key, while the API hands you `betId` on every transaction. Should these be two tables at all, or one ledger that expands?

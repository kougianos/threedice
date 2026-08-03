---
name: ThreeDice
description: A quiet casino terminal - dashboard sobriety in Midnight Slate, with one warm light on the dice.
colors:
  electric-iris: "#6c5ce7"
  electric-iris-lit: "#7f70f0"
  lucky-gold: "#ffd32a"
  winners-teal: "#00b894"
  ember-clay: "#e17055"
  iris-text: "#8f83f4"
  clay-text: "#e8846b"
  midnight-slate: "#0f1117"
  slate-surface: "#1a1d27"
  slate-raised: "#22263a"
  slate-line: "#2a2e3e"
  pearl: "#e4e6eb"
  fog: "#8b8fa3"
  fog-bright: "#a8acc0"
  on-iris: "#ffffff"
typography:
  display:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "2rem"
    fontWeight: 800
  numeral:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "1.8rem"
    fontWeight: 700
  headline:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "1.5rem"
    fontWeight: 700
  verdict:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "1.3rem"
    fontWeight: 700
  headline-compact:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "1.25rem"
    fontWeight: 700
  title:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "1.1rem"
    fontWeight: 600
  body:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "1rem"
    fontWeight: 400
    lineHeight: 1.6
  control:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "0.95rem"
    fontWeight: 500
  cell:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "0.9rem"
    fontWeight: 400
  label:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "0.85rem"
    fontWeight: 400
  label-caps:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "0.8rem"
    fontWeight: 500
    letterSpacing: "0.5px"
  caption:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "0.75rem"
    fontWeight: 400
  micro:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "0.7rem"
    fontWeight: 400
rounded:
  xs: "4px"
  sm: "6px"
  md: "8px"
  lg: "12px"
  pill: "20px"
spacing:
  xs: "0.4rem"
  sm: "0.6rem"
  control: "0.7rem"
  md: "1rem"
  lg: "1.5rem"
  xl: "2rem"
components:
  button-primary:
    backgroundColor: "{colors.electric-iris}"
    textColor: "{colors.on-iris}"
    rounded: "{rounded.md}"
    padding: "0.75rem 2rem"
    width: "100%"
  button-primary-hover:
    backgroundColor: "{colors.electric-iris-lit}"
    textColor: "{colors.on-iris}"
  button-quick-pick:
    backgroundColor: "{colors.midnight-slate}"
    textColor: "{colors.pearl}"
    rounded: "{rounded.md}"
    padding: "0.4rem 0.7rem"
  button-quick-pick-active:
    backgroundColor: "rgba(108, 92, 231, 0.2)"
    textColor: "{colors.on-iris}"
    rounded: "{rounded.md}"
  button-quick-pick-best:
    backgroundColor: "rgba(255, 211, 42, 0.08)"
    textColor: "{colors.pearl}"
    rounded: "{rounded.md}"
  button-quick-stake:
    backgroundColor: "{colors.midnight-slate}"
    textColor: "{colors.winners-teal}"
    rounded: "{rounded.sm}"
    padding: "0.35rem 0.7rem"
  input-field:
    backgroundColor: "{colors.midnight-slate}"
    textColor: "{colors.pearl}"
    rounded: "{rounded.md}"
    padding: "0.7rem 1rem"
    width: "100%"
  card:
    backgroundColor: "{colors.slate-surface}"
    textColor: "{colors.pearl}"
    rounded: "{rounded.lg}"
    padding: "{spacing.lg}"
  die-tile:
    backgroundColor: "{colors.midnight-slate}"
    textColor: "{colors.lucky-gold}"
    rounded: "{rounded.lg}"
    height: "64px"
    width: "64px"
  balance-pill:
    backgroundColor: "{colors.electric-iris}"
    textColor: "{colors.on-iris}"
    rounded: "{rounded.pill}"
    padding: "0.4rem 1rem"
  badge-won:
    backgroundColor: "rgba(0, 184, 148, 0.15)"
    textColor: "{colors.winners-teal}"
    rounded: "{rounded.lg}"
    padding: "0.15rem 0.6rem"
  badge-lost:
    backgroundColor: "rgba(225, 112, 85, 0.15)"
    textColor: "{colors.clay-text}"
    rounded: "{rounded.lg}"
    padding: "0.15rem 0.6rem"
  tab-active:
    textColor: "{colors.iris-text}"
  toast-error:
    backgroundColor: "{colors.ember-clay}"
    textColor: "{colors.on-iris}"
    rounded: "{rounded.md}"
    padding: "0.8rem 1.5rem"
---

# Design System: ThreeDice

## Overview

**Creative North Star: "The Quiet Casino Terminal"**

The room is dark and the instruments are on. Surfaces are Midnight Slate and its two lighter siblings, hairline
borders separate everything, and the type is the operating system's own UI face at rest. Then three dice land in
Lucky Gold and the terminal admits, briefly, that it is a casino. That contrast is the whole system: a composed
readout that keeps one warm light burning over the table.

It is precise and composed first. Values are exact, the rhythm repeats, controls answer in two tenths of a second
with a single pixel of lift, and nothing is approximate. But it is not solemn. The wordmark carries a drawn die,
the dice rest at `? × ? × ? = ?` before you have bet anything, and the quick picks tell you your real chance of
winning and what it pays before you spend a cent. The system knows it is a dice game and is willing to say so, as
long as it says so in the same measured voice it uses for the ledger. The wit is in the arithmetic, not in
borrowed glyphs: no emoji stands in for an icon, and no star stands in for a number.

Depth is tonal, never dramatic. The palette is disciplined about meaning: violet is the only color that acts, gold
is the only color that rolls, and teal and clay only ever report an outcome. Confirmed rejections: this system does
not become neon slot-machine chrome (gradient bevels, glow spam, confetti, coin shine), and it does not become a
crypto trading terminal (tickers, sparklines, wall-to-wall red and green, manufactured urgency).

**Key Characteristics:**

- One 900px column of stacked cards on a Midnight Slate ground, with a single width breakpoint at 600px and one
  input breakpoint at `pointer: coarse`.
- A three-step tonal ramp does all structural depth work; the one card shadow is atmosphere.
- Controls are punched into the card, not raised off it. Only the primary action rises.
- Gold has exactly four homes: three in the game, one in the identity that quotes it.
- Zero font requests. Zero image requests. The only artwork is one drawn 5-pip die, used as both favicon and mark.
- Every state change is 0.2s, with a 1px lift and a colored glow on the primary action.
- The equation is never hidden. It rests at `? × ? × ? = ?` and carries the prediction under the product, so the
  comparison always has both sides.
- Every number the interface states about odds is counted at runtime, never transcribed.

## Colors

A cool, low-light neutral field with four saturated colors, each holding a single job and never borrowing another's.

### Primary

- **Electric Iris** (`#6c5ce7`): the only color that means *action or identity*. It carries the primary button,
  the balance pill, the active tab underline, the selected quick pick, the mark in the wordmark, and the start of
  the wordmark gradient. **Iris Lit** (`#7f70f0`) is its hover state and appears nowhere else.
- **Iris Text** (`#8f83f4`): the same voice, lifted until it clears 4.5:1 as text on a card. Used wherever Iris has
  to be read rather than pressed: the active tab label, footer links, and the focus outline, which needs to stay
  legible against both the card and the page ground. Full-strength Iris stays for fills.

### Secondary

- **Lucky Gold** (`#ffd32a`): the game's own color. It appears on the three die numerals, on the product they
  multiply to, on the best-odds pick (as a border plus an 8% tint), and at the far end of the wordmark gradient.
  It is never a filled button, never body text, never a surface.

### Tertiary

The outcome pair. Both are strictly semantic and always appear as a 15% tint behind their own full-strength text.

- **Winner's Teal** (`#00b894`): won, credit, winnings, and the quick-stake buttons that add money to the field.
- **Ember Clay** (`#e17055`): lost, debit, borders on invalid fields, and the error toast, which is the one place
  clay fills a surface solid. **Clay Text** (`#e8846b`) is the lifted variant used whenever clay is read as text on
  its own tint, where full-strength clay only reaches 4.35:1.

### Neutral

- **Midnight Slate** (`#0f1117`): the page ground, and the recessed fill of every input, quick pick, and die tile.
- **Slate Surface** (`#1a1d27`): every card and the sticky header. The plane the interface actually lives on.
- **Slate Raised** (`#22263a`): hover fill for rows and picks, and the resting fill of odds badges.
- **Slate Line** (`#2a2e3e`): every border and divider, at 1px. Table row dividers use it at 50% opacity.
- **Pearl** (`#e4e6eb`): primary text.
- **Fog** (`#8b8fa3`): labels, table headers, muted metadata, the multiply and equals signs, the footer.
- **Fog Bright** (`#a8acc0`): the same role wherever the text sits on a color tint rather than a neutral surface.
  Plain Fog drops to 3.17:1 on the gold-tinted best pick, so the probability sublabel uses this instead.

### Named Rules

**The Iris Means Action Rule.** Violet marks what the visitor can do and who they are. If an element is not
actionable, not focused, and not the balance, it does not get Iris.

**The Gold Only Rolls Rule.** Lucky Gold has four homes and no more: the numerals of a roll, the product they make,
the picks marked best, and the far end of the wordmark gradient. The first three are the game itself; the fourth is
the identity quoting it. Adding a fifth costs the rest their meaning.

**The Outcome Pair Rule.** Teal and Clay report results and nothing else. They never become brand colors, never
decorate, and never appear on an element whose meaning is not won/lost or credit/debit.

## Typography

**Display Font:** none. The system uses one family throughout.
**Body Font:** `'Segoe UI', system-ui, -apple-system, sans-serif`
**Label/Mono Font:** none distinct. Labels are the same family, smaller, uppercased and tracked.

**Character:** the operating system's own UI voice, which reads as neutral and unbranded on purpose. Personality is
carried by size, weight and color rather than by letterforms: the family spans 400 to 800, and the jump from 400
body to 800 product is the loudest typographic move in the system. Every number that can be compared to another
number is set with `font-variant-numeric: tabular-nums`: money, balances, dice, products, probabilities and stakes.

### Hierarchy

Thirteen steps, each earning its place. The frontmatter is normative; this is what each one is for.

- **Display** (800, 2rem): the product value of a roll, in Lucky Gold. The largest and heaviest type in the system.
- **Numeral** (700, 1.8rem): the three die faces, in the same gold, one step below the product they make.
- **Headline** (700, 1.5rem): the wordmark, filled with a 135 degree Iris-to-Gold gradient clipped to the text.
  **Headline Compact** (700, 1.25rem) is the same wordmark under 600px.
- **Verdict** (700, 1.3rem): the won/lost pill. The multiply and equals operators share the size at 600 in Fog.
- **Title** (600, 1.1rem): card headings. Small enough that a card reads as a panel rather than a page section.
- **Body** (400, 1rem, line-height 1.6): the default, and the input value size, so a typed stake matches the
  surrounding text exactly. Also the primary button label at 600.
- **Control** (500, 0.95rem): tab labels.
- **Cell** (400, 0.9rem): table cells, the player name, and the value inside a quick pick at 600.
- **Label** (400, 0.85rem): field labels and the odds legend. Sentence case, not uppercased.
- **Label Caps** (500, 0.8rem, 0.5px tracking, uppercase): table column headers. Quick stakes, badges and field
  errors share the size at 600 without the uppercasing.
- **Caption** (400, 0.75rem): the prediction echo under the product, and the row labels the history table grows
  when it collapses to records on mobile.
- **Micro** (400, 0.7rem): the probability and payout line inside a quick pick. The floor; nothing goes smaller.

### Named Rules

**The Zero-Request Type Rule.** The system loads no font files. Type costs nothing at first paint today, and a
webfont is a trade to argue for rather than a default to reach for.

**The Weight Before Size Rule.** Emphasis is bought with weight (600, 700, 800) before it is bought with size. Only
the roll result is allowed to be both.

## Layout

A single 900px column, centered, with 1rem side padding, sitting under a sticky full-bleed header that shares the
same 900px inner measure. Cards stack vertically with a 1.5rem gap and 1.5rem internal padding, in a fixed order:
bet, result (hidden until the first roll), history. The body is a flex column with the footer pushed to the bottom,
so a short page still fills the viewport.

The only nested grid is the bet row: two equal columns (stake, prediction) at 1rem gap. Within each column, quick
buttons sit 0.4rem below their input in a wrapping flex row with 0.4rem gaps.

Spacing rhythm runs in 0.4rem / 0.6rem / 0.7rem / 1rem / 1.5rem / 2rem steps. Controls take 0.7rem vertical
padding, tables take 0.6rem cell padding, cards and the footer take 1.5rem, the header takes 1rem by 2rem.

One width breakpoint, at 600px. The header stays a single row and tightens to 0.6rem by 1rem, dropping the player
name so the balance keeps its place. The bet row collapses to one column, cards tighten to 1.25rem by 1rem, dice
shrink from 64px to 52px, and inputs and the primary button take a 48px minimum height.

**Tables become records below 600px.** Seven columns cannot survive a 375px viewport, and the table's natural width
used to push the whole document 63px wider than the screen. Each row becomes a bordered 8px card, the header row is
removed from view but kept for assistive tech, and every cell grows its own label from `data-label` on the left with
its value on the right. Above 600px the table stays a table, inside an `overflow-x: auto` wrapper so it can never
push the page sideways again.

One input breakpoint, independent of width: under `pointer: coarse` every quick stake, quick pick and tab takes a
44px minimum height. Touch gets the larger target whatever the screen size, and a mouse is never penalised for it.

### Named Rules

**The Single Column Rule.** Everything below the header is one column of cards at one measure. The bet row is the
only place the system nests a grid, and it collapses at the single breakpoint.

## Elevation & Depth

Tonal first. Depth is carried by the neutral ramp, not by shadow: Midnight Slate recedes, Slate Surface is the
resting plane, Slate Raised comes forward on hover. Because inputs and quick picks are filled with the page ground,
controls read as punched into the card rather than floating above it, and the primary button is the only element
that visibly sits on top of anything.

Shadow is atmosphere and never rank. Every card carries the same ambient shadow regardless of importance, and no
card is allowed a heavier one to look more significant.

### Shadow Vocabulary

- **Card ambient** (`box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3)`): on every card, always the same. Separates the
  surface from the ground without implying hierarchy.
- **Action glow** (`box-shadow: 0 4px 12px rgba(108, 92, 231, 0.4)`): primary button hover only, paired with a 1px
  lift. The only colored shadow in the system.
- **Focus outline** (`outline: 2px solid` Iris Text, `outline-offset: 2px`): every focusable element, without
  exception. Inputs additionally take a soft `0 0 0 3px rgba(108, 92, 231, 0.2)` ring and an Iris border, but that
  ring is atmosphere at 1.23:1 and never the indicator; the outline is.
- **Toast** (`box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4)`): the one floating element, slightly heavier than a card
  because it genuinely floats.

### Named Rules

**The Tone Ramp Rule.** Structural depth comes from Midnight, Surface and Raised. Shadow never ranks importance and
never gets heavier to make something matter more.

**The Inset Control Rule.** Anything that accepts input recedes to Midnight Slate with a Slate Line border. Only
the primary action rises off the surface, and it rises in Iris.

## Shapes

A five-step radius ladder, soft but never pill-shaped except where an element is deliberately a token: 4px on odds
badges, 6px on quick-stake buttons, 8px on inputs, buttons and quick picks, 12px on cards, dice and outcome badges,
20px on the balance and the result status.

Borders are the primary separator: 1px Slate Line on cards, inputs, quick picks, the header underline, table
headers and the footer rule, with row dividers at 50% opacity. Dice are the exception at 2px, which is what makes a
64px square read as an object rather than a container. No element uses a dashed, dotted or double border.

The favicon is the system's only mark: a 12px-radius violet square carrying five white pips in the die-face-of-five
arrangement. It shares the card radius exactly.

### Named Rules

**The Shared Corner Rule.** Cards and dice both take 12px, so a rolled die reads as a small card rather than a
foreign object dropped into the layout.

**The Pill Means Status Rule.** The 20px pill is reserved for things that report a live value: the balance and the
won/lost verdict. Actions never wear it.

## Components

### Buttons

- **Shape:** softly rounded (8px) for the primary action and quick picks; slightly tighter (6px) for quick stakes.
- **Primary:** full-width, Iris fill, white label at 600 / 1rem, 0.75rem by 2rem padding. One per card.
- **Hover / Focus:** background lifts to Iris Lit, the button rises 1px, and the action glow appears, all over 0.2s.
- **Disabled:** 50% opacity with a not-allowed cursor. The label also swaps to "Rolling..." while a bet is in
  flight, so the disabled state is never silent.
- **Quick stake:** Midnight fill, Slate Line border, Teal label at 600 / 0.8rem, reading `+$10`. On hover the
  border turns Teal, the fill takes a 10% Teal wash, and it lifts 1px. These add to the stake rather than set it.
- **Quick pick:** Midnight fill, Slate Line border, a two-line stack of the value at 600 / 0.9rem over its real win
  probability and payout at 400 / 0.7rem in Fog Bright, minimum 62px wide. Hover borders Iris and fills Slate
  Raised. Selection is `aria-pressed`, styled as an Iris border with a 20% Iris fill and white text. The `best`
  variant carries a Gold border with an 8% Gold fill at rest and 20% when selected.

  Both numbers on a pick are computed at load from the 216 ordered rolls, never transcribed. The set of picks
  marked `best` is whichever ties for the highest expected return, which is currently two of them, not one. The
  accessible name states the value first ("Predict 12. 6.9 percent chance, pays 5x"), because a pick that announces
  only its odds never says what it selects.

### Cards / Containers

- **Corner Style:** 12px.
- **Background:** Slate Surface on the Midnight ground.
- **Shadow Strategy:** card ambient, identical on every card. See Elevation & Depth.
- **Border:** 1px Slate Line, which does most of the separation work; the shadow only softens the edge.
- **Internal Padding:** 1.5rem, with a 1rem gap under the card heading.

### Inputs / Fields

- **Style:** full-width, Midnight fill, 1px Slate Line border, 8px radius, 0.7rem by 1rem padding, Pearl text at
  1rem with tabular figures. The label sits above in Pearl at 0.85rem and carries its constraint inline in Fog
  ("Stake $1.00 to $10,000.00", "Predicted product, one of 40 possible"), so the rule is readable before submitting
  rather than after failing.
- **Focus:** the Iris Text outline, plus an Iris border and the soft 3px ring, over 0.2s.
- **Error:** the field takes `aria-invalid` and a Clay border, and a Clay Text message appears directly beneath it
  at 0.8rem. Validation answers at the field it belongs to; the toast is for server and network failures only, and
  the form carries `novalidate` so the browser's own bubble is not a third channel.

### Navigation

The header is the only navigation: sticky, Slate Surface, 1px Slate Line underline, holding the gradient wordmark
on the left and the player identity on the right. The player name is Fog at 0.9rem; the balance is an Iris pill in
white 700. On a balance change the pill scales to 1.1 for 300ms and settles, which is the system's way of
confirming money moved.

The history card carries a second, local navigation under its own visible heading: text tabs at 500 / 0.95rem in
Fog on a Slate Line rule, with the active tab turning Iris Text and growing a 2px Iris underline. No pill, no fill,
no card. It is a real tab pattern, not two buttons that look like one: `role="tablist"` with `aria-selected`,
`aria-controls`, arrow-key movement, and a roving tabindex so the group is one tab stop rather than two.

### Tables

Column headers are Fog labels, uppercased, tracked 0.5px, on a Slate Line rule. Cells run 0.9rem with 0.6rem
padding and a 50%-opacity divider. Rows fill Slate Raised on hover. Outcome and transaction type are the only
colored cells, always as a 12px badge with a 15% tint. Empty states replace the whole row with centered italic Fog
text that names the next action ("No bets yet. Place your first bet above.").

Column order follows the comparison the row exists to make: **Bet, Predicted, Product, Dice, Stake, Result, Time**,
so the predicted value and the product it is being judged against sit side by side. The identity column carries the
server's real `betId` or `transactionId`, never a position in the list, and its header names which one it is; a
row's label must not change because a newer row arrived. Numeric columns are right-aligned with tabular figures.
Times are relative ("2m ago") with the absolute timestamp in `title`.

### Dice Equation (signature)

The result card is a single centered horizontal line that reads as arithmetic: three 64px Midnight tiles with 2px
Slate Line borders and Gold numerals at 700 / 1.8rem, separated by Fog multiply signs at 1.3rem, then a Fog equals
sign, then the product in Gold at 800 / 2rem. It wraps rather than scrolls, and at 600px the tiles drop to 52px.

**The card is never hidden.** It rests at `? × ? × ? = ?` from first paint, so the product's one authored idea is in
the opening viewport instead of behind a committed bet. Directly under the product sits the other half of the
comparison, a 0.75rem Fog caption reading "you predicted 12", which tracks the prediction field live and turns Teal
at 600 when the roll matches it. An equation is not an answer with one side missing.

The state machine has three positions, and the whole point is that they never overlap:

1. **Resting.** Dice and product at `?`, no verdict, a Fog line naming what happens next.
2. **Pending.** On submit the previous roll is cleared back to `?` and the row drops to 55% opacity. A bet in flight
   must never show the last bet's dice, product or verdict.
3. **Settled.** One 0.4s flip on all three tiles, a full 360 degree Y rotation with the scale dipping to 0.8 at the
   halfway point. Dice, product, prediction caption, verdict pill and winnings all land **together** at 200ms, mid
   flip. The verdict is never written before the dice that justify it.

Below sits the verdict: a 20px status pill (Teal on a 15% Teal tint with a matching border for a win, Clay Text on
a 15% Clay tint for a loss) beside the winnings in Teal at 600. Every settle also writes one sentence to the live
region: what was rolled, what it made, what was predicted, what it paid, and the new balance.

### Toast

The only floating element. Fixed 2rem from the bottom, centered, Ember Clay fill, white 500 / 0.9rem, 8px radius,
its own heavier shadow, capped at 30rem wide. It slides up from 100px below over 0.3s, holds for 8 seconds, and
leaves the same way. It carries `role="alert"` and a dismiss control, so it can be read at any pace and closed on
demand rather than expiring on a timer.

It is for **server and network failures only**. Anything the client can decide for itself is answered at the field.
Server detail strings are stripped of their `field: ` prefix before display, because an API key name is not a
sentence. The entrance class is applied after a forced reflow rather than inside `requestAnimationFrame`, which
never fires in a background tab and used to leave every toast parked below the viewport for its entire life.

### Named Rules

**The Two-Tenths Rule.** Every state change runs 0.2s. Only events are allowed longer: the roll at 0.4s, the toast
at 0.3s, the balance mark at 400ms. If it responds to a pointer, it is 0.2s.

**The One Loud Thing Rule.** Each card gets at most one element at full saturation: the Iris button in the bet
card, the Gold product in the result card, the badges in the history. Everything else in that card is neutral.
The odds legend is the standing test: it sits above the button in neutral badges precisely because a second
solid-Iris element in that card read as a second submit.

**The Nothing Precedes The Dice Rule.** No part of an outcome may appear before the dice that produce it. The
verdict, the winnings, the matched caption and the balance all wait for the same frame. A result on screen always
belongs to the bet the visitor is currently looking at, or the card is showing `?`.

**The Numbers Are Counted Rule.** Any probability, payout or count the interface states is derived at runtime from
the 216 ordered rolls, never transcribed into markup. Four of the eight quick picks once carried hand-written
percentages and four of them were wrong, one by a factor of two. A number the code cannot recompute is a number the
interface is not entitled to claim.

**The Direction Is The Feedback Rule.** The balance marks which way it moved: Teal and a 1.1 scale on a credit, a
0.97 settle on a debit. One gesture for both is a celebration attached to a loss.

## Do's and Don'ts

### Do:

- **Do** keep Lucky Gold on the game. Its four homes are the die numerals, the product, the best-odds picks, and
  the far end of the wordmark gradient.
- **Do** recess controls into Midnight Slate with a 1px Slate Line border, and let only the primary action rise.
- **Do** answer every pointer in 0.2s with a 1px lift. Crisp, with a spring in it.
- **Do** hold the 900px single column and let the one 600px breakpoint collapse the bet row and shrink the dice.
- **Do** reserve Teal and Clay for won/lost and credit/debit, always as a 15% tint under full-strength text.
- **Do** pair every quick pick with its real win probability and payout, both counted at runtime. Telling the
  visitor their odds is part of the voice; telling them wrongly is worse than silence.
- **Do** keep type on the system stack unless a webfont earns the request it costs.
- **Do** answer validation at the field that caused it, and keep the toast for failures the client cannot predict.
- **Do** give every focusable element the same Iris Text outline at 2px with a 2px offset.
- **Do** announce the outcome of a roll in the live region: rolled, product, predicted, paid, new balance.
- **Do** set every comparable number in tabular figures, right-aligned in tables.
- **Do** state a field's constraint in its label, before the visitor can break it.

### Don't:

- **Don't** drift into neon slot-machine chrome: no gradient bevels, glow spam, confetti, sparkles, or coin shine
  on the gold.
- **Don't** drift into a crypto trading terminal: no tickers, no sparklines, no wall-to-wall red and green
  numerals, no manufactured urgency.
- **Don't** rank importance with shadow. Every card carries the same ambient shadow; depth is tonal.
- **Don't** give Iris a second job. It means action, focus and identity, never decoration.
- **Don't** put two full-saturation elements in one card.
- **Don't** introduce a second colored shadow. The Iris action glow is the only one.
- **Don't** use the 20px pill for anything actionable; it belongs to live values.
- **Don't** write a probability, payout or count into markup by hand. Count it.
- **Don't** show any part of a result before the dice land, or leave a previous result on screen during a roll.
- **Don't** let a table push the document sideways. Wrap it, or collapse it to records.
- **Don't** identify a row by its position in the list. Use the id the service returned.
- **Don't** put a raw API field name in front of a visitor.
- **Don't** kill motion wholesale under `prefers-reduced-motion`. Remove the movement, keep the state change: the
  dice still turn over, the balance still marks its direction in color, the toast still arrives.

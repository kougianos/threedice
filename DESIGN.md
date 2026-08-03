---
name: ThreeDice
description: A quiet casino terminal - dashboard sobriety in Midnight Slate, with one warm light on the dice.
colors:
  electric-iris: "#6c5ce7"
  electric-iris-lit: "#7f70f0"
  lucky-gold: "#ffd32a"
  winners-teal: "#00b894"
  ember-clay: "#e17055"
  midnight-slate: "#0f1117"
  slate-surface: "#1a1d27"
  slate-raised: "#22263a"
  slate-line: "#2a2e3e"
  pearl: "#e4e6eb"
  fog: "#8b8fa3"
  on-iris: "#ffffff"
typography:
  display:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "2rem"
    fontWeight: 800
  headline:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "1.5rem"
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
  label:
    fontFamily: "'Segoe UI', system-ui, -apple-system, sans-serif"
    fontSize: "0.8rem"
    fontWeight: 500
    letterSpacing: "0.5px"
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
    padding: "0.3rem 0.6rem"
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
    textColor: "{colors.ember-clay}"
    rounded: "{rounded.lg}"
    padding: "0.15rem 0.6rem"
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
with a single pixel of lift, and nothing is approximate. But it is not solemn. The wordmark carries a die, results
land with an emoji, the best-odds pick wears a star, and the quick picks tell you your real chance of winning
before you spend anything. The system knows it is a dice game and is willing to say so, as long as it says so in
the same measured voice it uses for the ledger.

Depth is tonal, never dramatic. The palette is disciplined about meaning: violet is the only color that acts, gold
is the only color that rolls, and teal and clay only ever report an outcome. Confirmed rejections: this system does
not become neon slot-machine chrome (gradient bevels, glow spam, confetti, coin shine), and it does not become a
crypto trading terminal (tickers, sparklines, wall-to-wall red and green, manufactured urgency).

**Key Characteristics:**

- One 900px column of stacked cards on a Midnight Slate ground, with a single breakpoint at 600px.
- A three-step tonal ramp does all structural depth work; the one card shadow is atmosphere.
- Controls are punched into the card, not raised off it. Only the primary action rises.
- Gold has exactly three homes, all of them the game itself.
- Zero font requests. Zero image requests. The only artwork is a 5-pip SVG favicon.
- Every state change is 0.2s, with a 1px lift and a colored glow on the primary action.

## Colors

A cool, low-light neutral field with four saturated colors, each holding a single job and never borrowing another's.

### Primary

- **Electric Iris** (`#6c5ce7`): the only color that means *action or identity*. It carries the primary button,
  the balance pill, the active tab underline, the focus ring, the 5x odds highlight, the selected quick pick, and
  the start of the wordmark gradient. **Iris Lit** (`#7f70f0`) is its hover state and appears nowhere else.

### Secondary

- **Lucky Gold** (`#ffd32a`): the game's own color. It appears on the three die numerals, on the product they
  multiply to, on the best-odds pick (as a border plus an 8% tint), and at the far end of the wordmark gradient.
  It is never a filled button, never body text, never a surface.

### Tertiary

The outcome pair. Both are strictly semantic and always appear as a 15% tint behind their own full-strength text.

- **Winner's Teal** (`#00b894`): won, credit, winnings, and the quick-stake buttons that add money to the field.
- **Ember Clay** (`#e17055`): lost, debit, and the error toast, which is the one place clay fills a surface solid.

### Neutral

- **Midnight Slate** (`#0f1117`): the page ground, and the recessed fill of every input, quick pick, and die tile.
- **Slate Surface** (`#1a1d27`): every card and the sticky header. The plane the interface actually lives on.
- **Slate Raised** (`#22263a`): hover fill for rows and picks, and the resting fill of odds badges.
- **Slate Line** (`#2a2e3e`): every border and divider, at 1px. Table row dividers use it at 50% opacity.
- **Pearl** (`#e4e6eb`): primary text.
- **Fog** (`#8b8fa3`): labels, table headers, muted metadata, the multiply and equals signs, the footer.

### Named Rules

**The Iris Means Action Rule.** Violet marks what the visitor can do and who they are. If an element is not
actionable, not focused, and not the balance, it does not get Iris.

**The Gold Only Rolls Rule.** Lucky Gold has three homes: the numerals of a roll, the product they make, and the
pick marked best. Adding a fourth costs the first three their meaning.

**The Outcome Pair Rule.** Teal and Clay report results and nothing else. They never become brand colors, never
decorate, and never appear on an element whose meaning is not won/lost or credit/debit.

## Typography

**Display Font:** none. The system uses one family throughout.
**Body Font:** `'Segoe UI', system-ui, -apple-system, sans-serif`
**Label/Mono Font:** none distinct. Labels are the same family, smaller, uppercased and tracked.

**Character:** the operating system's own UI voice, which reads as neutral and unbranded on purpose. Personality is
carried by size, weight and color rather than by letterforms: the family spans 400 to 800, and the jump from 400
body to 800 product is the loudest typographic move in the system. Numerals are proportional, and the history
tables do not currently request tabular figures.

### Hierarchy

- **Display** (800, 2rem): the product value of a roll. The largest and heaviest type in the interface, in Lucky
  Gold. The three die numerals sit just under it at 700 / 1.8rem in the same gold.
- **Headline** (700, 1.5rem): the wordmark only, filled with a 135 degree Iris-to-Gold gradient clipped to the text.
- **Title** (600, 1.1rem): card headings. Small enough that a card reads as a panel rather than a page section.
- **Body** (400, 1rem, line-height 1.6): the default. Also the input value size, so typed stake and prediction
  match the surrounding text weight exactly.
- **Label** (500, 0.8rem, 0.5px tracking, uppercase): table column headers. Field labels are the same size family
  at 0.85rem in Fog, sentence case, not uppercased.

Two supporting sizes recur: 0.9rem for table cells, tabs and secondary text, and 0.65rem for the probability
sublabel inside a quick pick.

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

One breakpoint, at 600px: the header stacks into two centered rows, the bet row collapses to one column, dice
shrink from 64px to 50px and their gap tightens, and table type drops to 0.8rem with 0.4rem cell padding.

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
- **Focus ring** (`box-shadow: 0 0 0 3px rgba(108, 92, 231, 0.2)`): inputs on focus, with the border shifting to
  Iris at the same time.
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
- **Quick pick:** Midnight fill, Slate Line border, a two-line stack of the value at 600 / 0.9rem over its win
  probability at 400 / 0.65rem in Fog, minimum 52px wide. Hover borders Iris and fills Slate Raised. Selected takes
  an Iris border with a 20% Iris fill and white text. The `best` variant carries a Gold border with an 8% Gold
  fill at rest and 20% when selected, and its sublabel ends in a star.

### Cards / Containers

- **Corner Style:** 12px.
- **Background:** Slate Surface on the Midnight ground.
- **Shadow Strategy:** card ambient, identical on every card. See Elevation & Depth.
- **Border:** 1px Slate Line, which does most of the separation work; the shadow only softens the edge.
- **Internal Padding:** 1.5rem, with a 1rem gap under the card heading.

### Inputs / Fields

- **Style:** full-width, Midnight fill, 1px Slate Line border, 8px radius, 0.7rem by 1rem padding, Pearl text at
  1rem. Labels sit above at 0.85rem in Fog.
- **Focus:** the outline is removed and replaced by an Iris border plus the 3px Iris focus ring, over 0.2s.
- **Error:** fields do not carry an error state. Validation failures surface in the toast instead.

### Navigation

The header is the only navigation: sticky, Slate Surface, 1px Slate Line underline, holding the gradient wordmark
on the left and the player identity on the right. The player name is Fog at 0.9rem; the balance is an Iris pill in
white 700. On a balance change the pill scales to 1.1 for 300ms and settles, which is the system's way of
confirming money moved.

The history card carries a second, local navigation: text tabs at 500 / 0.95rem in Fog on a Slate Line rule, with
the active tab turning Iris and growing a 2px Iris underline. No pill, no fill, no card.

### Tables

Column headers are Fog labels, uppercased, tracked 0.5px, on a Slate Line rule. Cells run 0.9rem with 0.6rem
padding and a 50%-opacity divider. Rows fill Slate Raised on hover. Outcome and transaction type are the only
colored cells, always as a 12px badge with a 15% tint. Empty states replace the whole row with centered italic Fog
text that names the next action ("No bets yet. Place your first bet above!").

### Dice Equation (signature)

The result card is a single centered horizontal line that reads as arithmetic: three 64px Midnight tiles with 2px
Slate Line borders and Gold numerals at 700 / 1.8rem, separated by Fog multiply signs at 1.3rem, then a Fog equals
sign, then the product in Gold at 800 / 2rem. It wraps rather than scrolls, and at 600px the tiles drop to 50px.

Below it sits the verdict: a 20px status pill (Teal on a 15% Teal tint with a matching border for a win, Clay for a
loss) beside the winnings in Teal at 600. Each roll replays a 0.4s animation on all three tiles: a full 360 degree
Y rotation with the scale dipping to 0.8 at the halfway point, with the values swapping in at 200ms so the numbers
change mid-flip rather than before it.

### Toast

The only floating element. Fixed 2rem from the bottom, centered, Ember Clay fill, white 500 / 0.9rem, 8px radius,
its own heavier shadow. It slides up from 100px below over 0.3s, holds for 3 seconds, slides back and removes
itself. It is the single channel for every error, including server validation detail strings.

### Named Rules

**The Two-Tenths Rule.** Every state change runs 0.2s. Only events are allowed longer: the roll at 0.4s, the toast
at 0.3s, the balance pop at 300ms. If it responds to a pointer, it is 0.2s.

**The One Loud Thing Rule.** Each card gets at most one element at full saturation: the Iris button in the bet
card, the Gold product in the result card, the badges in the history. Everything else in that card is neutral.

## Do's and Don'ts

### Do:

- **Do** keep Lucky Gold on the game. Its three homes are the die numerals, the product, and the best-odds pick.
- **Do** recess controls into Midnight Slate with a 1px Slate Line border, and let only the primary action rise.
- **Do** answer every pointer in 0.2s with a 1px lift. Crisp, with a spring in it.
- **Do** hold the 900px single column and let the one 600px breakpoint collapse the bet row and shrink the dice.
- **Do** reserve Teal and Clay for won/lost and credit/debit, always as a 15% tint under full-strength text.
- **Do** pair every quick pick with its real win probability. Telling the visitor their odds is part of the voice.
- **Do** keep type on the system stack unless a webfont earns the request it costs.

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

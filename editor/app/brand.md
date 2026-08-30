# Brand — LevelCraft

_Status: set_

Full research and rationale live in `docs/client/personality.md` — this
file is the quick-reference token summary for future work in this repo.

## Palette (dark-only; see personality.md for why)

| Token | Value | Use |
|---|---|---|
| `--background` | `oklch(0.145 0 0)` | App background |
| `--card` / `--popover` | `oklch(0.205 0 0)` | Elevated surfaces |
| `--secondary` / `--muted` / `--accent` (neutral) | `oklch(0.269 0 0)` | Hover/secondary surfaces |
| `--foreground` | `oklch(0.985 0 0)` | Primary text |
| `--muted-foreground` | `oklch(0.708 0 0)` | Secondary text |
| `--border` / `--input` | `oklch(1 0 0 / 10-15%)` | Hairlines |
| `--primary` (brand accent) | `oklch(0.74 0.14 55)` — warm amber | CTAs, active states |
| `--primary-foreground` | `oklch(0.145 0 0)` | Text on the accent (dark, not white) |
| `--ring` | matches `--primary` | Focus rings |
| `--destructive` | `oklch(0.704 0.191 22.216)` | Destructive actions only |

One accent color (`--primary`), used identically everywhere it appears —
never a second "brand" color for a different purpose.

## Typography

No custom webfont — Tailwind v4's default system stacks
(`font-sans` / `font-mono`), deliberately. See personality.md §Typography
for the reasoning (native rendering, offline-first, OS accessibility
settings).

## Voice

Direct, plain-sentence, one small human touch in empty/error states.
Never corporate ("Click here to..."), never cutesy/gamer-slang.

## Density

List-based over card-based (no scene thumbnails exist to put on a card
yet). Full-bleed screens, not centered small-card layouts.

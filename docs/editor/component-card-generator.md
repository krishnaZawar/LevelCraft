# Component Card Generator — Design

How the Attributes panel builds a component's card, and why it is driven by a
config file rather than by code per component.

For the config language itself, see
[component-card-config.md](./component-card-config.md). This document is about
the approach.

## The problem

Before this, every component's card was rendered by one hardcoded path in
`EditorShell.tsx`: a list of field names per component, and a single widget —
a number input — for all of them.

That held while the engine had Transform and Color, which are nothing but
numbers. It broke the moment `InputMovement` arrived. Its data is a mix of a
number, two booleans and four key names, so the panel drew seven number inputs,
showed `0` for each boolean, and offered no way to reach the key fields at all.
The component was in the engine and unusable from the editor.

The deeper problem is that the panel had to learn each component. Every new
component meant editing a shared view file, and the person adding it had to
understand how cards were assembled before they could add one. Three things
followed from that:

- **Cards drift.** Two people lay out two components two ways, and the panel
  slowly stops looking like one panel.
- **The bar to adding a component is the panel's internals**, not the
  component.
- **Reviewing a card is reviewing React**, which hides a layout question inside
  a code question.

## The approach

Cards are generated from `editor/app/config/component-cards.yaml`. The config
says which fields a component shows, how they are grouped into rows, and which
widget edits each one. The generator reads it once at startup and builds every
card through the same code path.

```
component-cards.yaml ──► validate ──► generator ──► widgets ──► card
                                          ▲
                        component data ───┘
```

Three consequences, which are the point of the change:

- **Consistency is structural.** There is one renderer, so a slider is the same
  slider everywhere. Nobody can lay a card out inconsistently without changing
  the generator for everyone.
- **Adding a component is a config edit.** No React, no shared view file, and
  the reviewer reads a layout rather than a component tree.
- **The codebase stops growing per component.** A new component adds rows to a
  config file, not a new branch in a view.

## What the generator owns, and what it does not

The engine remains authoritative over the data. A component's fields, their
types, and what counts as a valid value are decided in Go, and the backend
rejects anything invalid regardless of what the editor drew. The config decides
presentation only: grouping, widget, label, and whether a field is shown at
all.

This split matters because the alternative — letting the config describe the
data too — would make the YAML a second source of truth that has to agree with
Go. Where the config does restate an engine constraint, like Color's 0–255
bounds, it does so to shape the widget rather than to enforce the rule.

Option lists that the engine owns are not copied into the config. A field
declares *where* its options come from and the editor fetches them, so the key
mapping lives in the backend and the config only points at it.

## Validation

The config is validated when it loads, and a bad config fails loudly rather
than rendering a partial card. A card that silently drops a field is worse than
one that refuses to load, because the field looks absent rather than broken.

Every problem is collected and reported together, so a config with three
mistakes takes one fix round rather than three.

The checks that matter are the ones that catch a mistake the author cannot see
in the file: a slider with only one bound, a select with no options, a
condition on a field the card never shows, two inputs bound to the same key.

## Components with no config

A component with no entry renders anyway. The generator falls back to the
fields the component actually carries and picks a widget from each value's
type.

This is deliberate. It means adding a component to the engine never makes it
invisible in the editor — it just makes it unstyled until someone lays it out —
and it means the config is an improvement over the default rather than a
requirement before anything works. The same fallback applies to a single field
added to an already-configured component.

## Co-dependent fields

Some components validate several fields together. `InputMovement` rejects an
axis switched on unless that axis's keys are already set, so a card that
committed one field at a time could never enable it: the flag alone is invalid,
and the keys are not reachable until the flag is on.

The generator therefore commits a **patch** rather than a field, and a checkbox
can declare the values that must accompany it. This keeps the knowledge of
which fields travel together in the config, next to the layout that makes them
co-dependent, rather than in the component that happens to render them.

## Trade-offs

**A config file is not type-checked.** A typo in a field name is caught by
validation at startup rather than by the compiler. This is the cost of the
config being data; the validator exists to pay it down.

**The config can drift from the engine.** Nothing stops a field being removed
in Go while the config still lays it out. Today that row renders empty. Having
the backend serve each component's field schema would close this, and would
also remove the duplicated range constants — worth doing, and deliberately not
in this change.

**Not every card will fit.** A colour swatch or an asset picker will eventually
want to render itself. The generator should grow an escape hatch for a named
custom renderer before that happens, rather than growing the config language
until it can express every widget.

## Where things live

| Path | What it is |
| --- | --- |
| `config/component-cards.yaml` | The cards, as data |
| `src/renderer/src/cards/types.ts` | The config's shape |
| `src/renderer/src/cards/validate.ts` | Validation, with the rules in one place |
| `src/renderer/src/cards/config.ts` | Loading, and the fallback for unconfigured components |
| `src/renderer/src/cards/fields.tsx` | One widget per input type |
| `src/renderer/src/cards/CardBody.tsx` | The generator |

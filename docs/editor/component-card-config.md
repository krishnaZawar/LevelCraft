# Component Card Config — Language Reference

How to lay out a component's card in `editor/app/config/component-cards.yaml`.

For why cards are built this way, see
[component-card-generator.md](./component-card-generator.md).

## Adding a component

A component the engine ships but the config never mentions still renders: the
generator falls back to the fields it carries and picks a widget from each
value's type. So the first question is not "how do I make this appear" but
"how do I lay this out deliberately".

A minimal entry:

```yaml
version: 1

components:
  Health:
    rows:
      - fields:
          - { field: current, label: Current, input: number, min: 0 }
          - { field: max, label: Max, input: number, min: 1 }
```

`Health` must match the component's name in Go exactly, and each `field` must
match a key in that component's serialised data.

## Top level

| Key | Required | Notes |
| --- | --- | --- |
| `version` | yes | Currently `1`. A file declaring any other version is rejected rather than guessed at. |
| `components` | yes | Keyed by component name. |

### A component entry

| Key | Required | Notes |
| --- | --- | --- |
| `description` | no | Shown as a tooltip on the card heading. |
| `rows` | yes | Rendered top to bottom, in list order. |

## Rows

A row is a horizontal group of fields, holding **1 to 4** of them — past four
the panel is too narrow to stay legible, and the config is rejected.

```yaml
rows:
  - fields:
      - { field: x, label: X, input: number }
      - { field: y, label: Y, input: number }
  - fields:
      - { field: w, label: W, input: number }
      - { field: h, label: H, input: number }
```

Order is the list's order. There are no `Row 1` / `Row 2` keys, because naming
a row only to sort by that name makes `Row 10` sort before `Row 2`.

| Key | Required | Notes |
| --- | --- | --- |
| `fields` | yes | 1–4 fields. |
| `visibleIf` | no | Hides the whole row. See [Conditions](#conditions). |

## Fields

| Key | Required | Default | Notes |
| --- | --- | --- | --- |
| `field` | yes | — | The key in the component's data. |
| `input` | yes | — | One of the widgets below. |
| `label` | no | the field name | Shown beside the input. |
| `help` | no | — | Tooltip. |
| `visibleIf` | no | — | Hides this field only. |

### Widgets

| `input` | Edits | Extra keys |
| --- | --- | --- |
| `number` | a number | `min`, `max`, `step` |
| `slider` | a number | `min`, `max` (**both required**), `step` |
| `text` | a string | — |
| `checkbox` | a boolean | `enables` |
| `select` | a string | `options` or `optionsFrom` (**one required**) |

**`number`** commits on blur and on Enter, reverts on Escape, and clamps to
`min`/`max` rather than rejecting — so a typed-in `900` becomes `255` instead
of bouncing off the backend.

**`slider`** needs both bounds because one bound is not a track. It commits
when the drag ends, not per pixel, so dragging sends one update rather than one
per value it passes through.

**`select`** takes either a fixed list:

```yaml
- field: axis
  label: Axis
  input: select
  options:
    - { value: x, label: Horizontal }
    - { value: y, label: Vertical }
```

`label` defaults to `value` when omitted.

…or a list the engine owns:

```yaml
- { field: leftKey, label: Left, input: select, optionsFrom: inputKeys }
```

`optionsFrom` is how a list stays out of this file. `inputKeys` is the engine's
key mapping, fetched from the backend, so adding a key to the engine does not
mean editing this config. It is currently the only source.

## Conditions

`visibleIf` shows a field or a row only when another field on the **same
component** has a given value:

```yaml
- visibleIf: { field: allowX, equals: true }
  fields:
    - { field: leftKey, label: Left, input: select, optionsFrom: inputKeys }
```

| Key | Notes |
| --- | --- |
| `field` | Another field on this component. It must be laid out somewhere on the card, or the config is rejected — a condition on a field the card never shows can never become true. |
| `equals` | A string, number or boolean. |

Deliberately one comparison rather than an expression language: it can be
checked when the config loads, and it cannot loop.

Hiding a field never deletes its value.

## `enables` — fields that must travel together

Some components validate a flag together with the fields it guards.
`InputMovement` rejects `allowX: true` unless `leftKey` and `rightKey` are
already valid, so switching the checkbox on by itself would always fail.

`enables` lets a checkbox carry those values in the same update:

```yaml
- field: allowX
  label: Move horizontally
  input: checkbox
  enables: { leftKey: KEY_A, rightKey: KEY_D }
```

Ticking the box commits `allowX: true` **and** the two keys at once. Unticking
commits only `allowX: false`, leaving the component to clear whatever it wants.

`enables` applies to a `checkbox` only, and every key it names must be laid out
on the same card.

## What is rejected

Validation runs once at startup and reports every problem together, so a config
with three mistakes takes one fix round. Any of these refuses to load:

- a `version` other than `1`
- a row with no fields, or more than four
- an unknown `input`
- the same `field` laid out twice on one component
- a `slider` missing `min` or `max`, or a `min` above its `max`
- a `select` with neither `options` nor `optionsFrom`
- an unknown `optionsFrom`
- `visibleIf` naming a field the card never shows
- `enables` on something that is not a checkbox, or naming a field the card never shows

A field the component carries but the config never mentions is **not** an
error. It is appended in its own row with a widget chosen from its type, so
adding a field in Go never makes it invisible — only unstyled until someone
lays it out.

## Worked example

```yaml
version: 1

components:
  Transform:
    description: Position and size in world space.
    rows:
      - fields:
          - { field: x, label: X, input: number }
          - { field: y, label: Y, input: number }
      - fields:
          - { field: w, label: W, input: number, min: 1 }
          - { field: h, label: H, input: number, min: 1 }

  Color:
    rows:
      - fields:
          - { field: r, label: R, input: slider, min: 0, max: 255 }
          - { field: g, label: G, input: slider, min: 0, max: 255 }
      - fields:
          - { field: b, label: B, input: slider, min: 0, max: 255 }
          - { field: a, label: A, input: slider, min: 0, max: 255 }

  InputMovement:
    description: Move the object with the keyboard.
    rows:
      - fields:
          - { field: speed, label: Speed, input: number, min: 0 }
      - fields:
          - field: allowX
            label: Move horizontally
            input: checkbox
            enables: { leftKey: KEY_A, rightKey: KEY_D }
      - visibleIf: { field: allowX, equals: true }
        fields:
          - { field: leftKey, label: Left, input: select, optionsFrom: inputKeys }
          - { field: rightKey, label: Right, input: select, optionsFrom: inputKeys }
```

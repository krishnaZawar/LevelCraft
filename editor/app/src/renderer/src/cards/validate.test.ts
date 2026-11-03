import { describe, expect, it } from 'vitest'
import { parse } from 'yaml'
import { type CardConfigError, validateCardConfig } from './validate'

function validate(yaml: string): ReturnType<typeof validateCardConfig> {
  return validateCardConfig(parse(yaml))
}

function problemsFrom(yaml: string): string {
  try {
    validate(yaml)
  } catch (err) {
    return (err as CardConfigError).message
  }
  throw new Error('expected the config to be rejected, but it passed')
}

const VALID = `
version: 1
components:
  Transform:
    rows:
      - fields:
          - { field: x, label: X, input: number }
          - { field: y, label: Y, input: number }
`

describe('a valid config', () => {
  it('keeps the rows in the order they were written', () => {
    const config = validate(`
version: 1
components:
  Transform:
    rows:
      - fields: [{ field: x, label: X, input: number }]
      - fields: [{ field: y, label: Y, input: number }]
`)
    const fields = config.components.Transform.rows.map((row) => row.fields[0].field)
    expect(fields).toEqual(['x', 'y'])
  })

  it('falls back to the field name when no label is given', () => {
    const config = validate(`
version: 1
components:
  Transform:
    rows:
      - fields: [{ field: speed, input: number }]
`)
    expect(config.components.Transform.rows[0].fields[0].label).toBe('speed')
  })
})

describe('structural rules', () => {
  it('rejects a version it does not understand', () => {
    expect(problemsFrom(VALID.replace('version: 1', 'version: 2'))).toContain('version must be 1')
  })

  it('rejects an unknown input type', () => {
    expect(problemsFrom(VALID.replace('input: number', 'input: wysiwyg'))).toContain(
      'unknown input "wysiwyg"'
    )
  })

  it('rejects a row with more fields than fit', () => {
    const tooWide = `
version: 1
components:
  Transform:
    rows:
      - fields:
          - { field: a, input: number }
          - { field: b, input: number }
          - { field: c, input: number }
          - { field: d, input: number }
          - { field: e, input: number }
`
    expect(problemsFrom(tooWide)).toContain('a row holds 1 to 4 fields')
  })

  // Two inputs bound to one key would fight over the same value.
  it('rejects the same field laid out twice', () => {
    const duplicated = `
version: 1
components:
  Transform:
    rows:
      - fields: [{ field: x, input: number }]
      - fields: [{ field: x, input: number }]
`
    expect(problemsFrom(duplicated)).toContain('"x" is laid out more than once')
  })

  // One fix round rather than three.
  it('reports every problem at once', () => {
    const broken = `
version: 9
components:
  Transform:
    rows:
      - fields: [{ field: x, input: nope }]
`
    const message = problemsFrom(broken)
    expect(message).toContain('version must be 1')
    expect(message).toContain('unknown input "nope"')
  })
})

describe('widget-specific rules', () => {
  it('requires both bounds on a slider, since one bound is not a track', () => {
    const halfBounded = `
version: 1
components:
  Color:
    rows:
      - fields: [{ field: r, input: slider, min: 0 }]
`
    expect(problemsFrom(halfBounded)).toContain('a slider needs both "min" and "max"')
  })

  it('rejects a min above its max', () => {
    const inverted = `
version: 1
components:
  Color:
    rows:
      - fields: [{ field: r, input: slider, min: 255, max: 0 }]
`
    expect(problemsFrom(inverted)).toContain('min (255) is above max (0)')
  })

  it('requires a select to say where its options come from', () => {
    const optionless = `
version: 1
components:
  InputMovement:
    rows:
      - fields: [{ field: leftKey, input: select }]
`
    expect(problemsFrom(optionless)).toContain('a select needs "options" or "optionsFrom"')
  })

  it('rejects an unknown dynamic option source', () => {
    const unknownSource = `
version: 1
components:
  InputMovement:
    rows:
      - fields: [{ field: leftKey, input: select, optionsFrom: gamepadButtons }]
`
    expect(problemsFrom(unknownSource)).toContain('unknown optionsFrom "gamepadButtons"')
  })

  it('defaults an option label to its value', () => {
    const config = validate(`
version: 1
components:
  InputMovement:
    rows:
      - fields:
          - field: axis
            input: select
            options: [{ value: x }, { value: y, label: Vertical }]
`)
    expect(config.components.InputMovement.rows[0].fields[0].options).toEqual([
      { value: 'x', label: 'x' },
      { value: 'y', label: 'Vertical' }
    ])
  })
})

describe('conditions and enables', () => {
  // A condition on a field the card never shows can never become true, which
  // makes the row it guards dead rather than conditional.
  it('rejects a condition on a field the card never shows', () => {
    const dangling = `
version: 1
components:
  InputMovement:
    rows:
      - visibleIf: { field: allowZ, equals: true }
        fields: [{ field: leftKey, input: text }]
`
    expect(problemsFrom(dangling)).toContain('visibleIf refers to "allowZ"')
  })

  it('accepts a condition on a field laid out later in the card', () => {
    expect(() =>
      validate(`
version: 1
components:
  InputMovement:
    rows:
      - visibleIf: { field: allowX, equals: true }
        fields: [{ field: leftKey, input: text }]
      - fields: [{ field: allowX, input: checkbox }]
`)
    ).not.toThrow()
  })

  it('rejects enables naming a field the card never shows', () => {
    const dangling = `
version: 1
components:
  InputMovement:
    rows:
      - fields:
          - { field: allowX, input: checkbox, enables: { leftKey: KEY_A } }
`
    expect(problemsFrom(dangling)).toContain('enables "leftKey"')
  })

  it('rejects enables on anything but a checkbox', () => {
    const misplaced = `
version: 1
components:
  InputMovement:
    rows:
      - fields:
          - { field: speed, input: number, enables: { leftKey: KEY_A } }
          - { field: leftKey, input: text }
`
    expect(problemsFrom(misplaced)).toContain('enables only applies to a checkbox')
  })
})

import { describe, expect, it } from 'vitest'
import { cardConfig, getCardConfig, resolveRows } from './config'

// Importing the module at all parses and validates the real shipped config, so
// a broken component-cards.yaml fails here rather than in front of a user.
describe('the shipped config', () => {
  it('loads and covers the components the engine ships', () => {
    expect(Object.keys(cardConfig.components).sort()).toEqual([
      'Color',
      'InputMovement',
      'Transform'
    ])
  })

  it('gates the movement keys behind their axis', () => {
    const rows = getCardConfig('InputMovement')?.rows ?? []
    const keyRow = rows.find((row) => row.fields.some((field) => field.field === 'leftKey'))
    expect(keyRow?.visibleIf).toEqual({ field: 'allowX', equals: true })
  })

  // The engine rejects an axis switched on without its keys, so the checkbox
  // has to carry them.
  it('names default keys on the axis checkboxes', () => {
    const rows = getCardConfig('InputMovement')?.rows ?? []
    const allowX = rows.flatMap((row) => row.fields).find((field) => field.field === 'allowX')
    expect(allowX?.enables).toEqual({ leftKey: 'KEY_A', rightKey: 'KEY_D' })
  })
})

describe('resolveRows', () => {
  it('lays out a configured component the way the config says', () => {
    const { rows } = resolveRows('Transform', { x: 1, y: 2, w: 3, h: 4 })
    expect(rows.map((row) => row.fields.map((field) => field.field))).toEqual([
      ['x', 'y'],
      ['w', 'h']
    ])
  })

  // A component nobody has laid out yet still has to be editable, or adding
  // one to the engine would make it invisible in the editor.
  it('falls back to the fields a component actually carries', () => {
    const { rows } = resolveRows('NotConfigured', { speed: 10, enabled: true, tag: 'hero' })
    expect(rows.map((row) => row.fields[0])).toEqual([
      expect.objectContaining({ field: 'speed', input: 'number' }),
      expect.objectContaining({ field: 'enabled', input: 'checkbox' }),
      expect.objectContaining({ field: 'tag', input: 'text' })
    ])
  })

  // Adding a field in Go should never make it invisible — just unstyled until
  // someone lays it out.
  it('appends fields the config does not mention', () => {
    const { rows } = resolveRows('Transform', { x: 1, y: 2, w: 3, h: 4, rotation: 90 })
    const laidOut = rows.flatMap((row) => row.fields.map((field) => field.field))
    expect(laidOut).toEqual(['x', 'y', 'w', 'h', 'rotation'])
  })

  it('does not invent rows for fields the component is missing', () => {
    const { rows } = resolveRows('Transform', { x: 1, y: 2, w: 3, h: 4 })
    expect(rows).toHaveLength(2)
  })
})

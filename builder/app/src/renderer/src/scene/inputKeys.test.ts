import { describe, expect, it } from 'vitest'
import { engineKeyName, engineMouseName } from './inputKeys'
import mapping from '../../../../../../utils/input/mapping.json'

const ENGINE_KEYS = new Set([...mapping.keyboard, ...mapping.mouse].map((key) => key.name))

describe('engineKeyName', () => {
  it('derives letters and digits from the code', () => {
    expect(engineKeyName('KeyA')).toBe('KEY_A')
    expect(engineKeyName('KeyZ')).toBe('KEY_Z')
    expect(engineKeyName('Digit0')).toBe('KEY_0')
    expect(engineKeyName('Digit9')).toBe('KEY_9')
  })

  it('derives function keys', () => {
    expect(engineKeyName('F1')).toBe('KEY_F1')
    expect(engineKeyName('F12')).toBe('KEY_F12')
  })

  // The engine's codes are GLFW's, so arrows cannot be matched numerically —
  // this is exactly where a browser keyCode would have gone wrong.
  it('maps arrows to the engine names rather than their codes', () => {
    expect(engineKeyName('ArrowLeft')).toBe('KEY_LEFT')
    expect(engineKeyName('ArrowUp')).toBe('KEY_UP')
  })

  it('maps the side-specific modifiers', () => {
    expect(engineKeyName('ShiftLeft')).toBe('KEY_LEFT_SHIFT')
    expect(engineKeyName('ControlRight')).toBe('KEY_RIGHT_CONTROL')
    expect(engineKeyName('MetaLeft')).toBe('KEY_LEFT_SUPER')
  })

  it('treats the numpad Enter as Enter', () => {
    expect(engineKeyName('NumpadEnter')).toBe('KEY_ENTER')
  })

  it('returns null for a key the engine has no name for', () => {
    expect(engineKeyName('CapsLock')).toBeNull()
    expect(engineKeyName('F13')).toBeNull()
    expect(engineKeyName('Numpad5')).toBeNull()
  })

  // The point of deriving rather than hand-copying: every name it produces has
  // to be one the engine actually ships.
  it('only ever produces names the engine knows', () => {
    const codes = [
      ...Array.from({ length: 26 }, (_, i) => `Key${String.fromCharCode(65 + i)}`),
      ...Array.from({ length: 10 }, (_, i) => `Digit${i}`),
      ...Array.from({ length: 12 }, (_, i) => `F${i + 1}`),
      'Space',
      'Escape',
      'Enter',
      'Tab',
      'Backspace',
      'Insert',
      'Delete',
      'ArrowRight',
      'ArrowLeft',
      'ArrowDown',
      'ArrowUp',
      'PageUp',
      'PageDown',
      'Home',
      'End',
      'ShiftLeft',
      'ControlLeft',
      'AltLeft',
      'MetaLeft',
      'ShiftRight',
      'ControlRight',
      'AltRight',
      'MetaRight'
    ]

    const produced = codes.map(engineKeyName).filter((name): name is string => name !== null)
    expect(produced.filter((name) => !ENGINE_KEYS.has(name))).toEqual([])
    // and it reaches every keyboard key the engine ships
    expect(produced.length).toBe(mapping.keyboard.length)
  })
})

describe('engineMouseName', () => {
  it('maps the DOM button numbers', () => {
    expect(engineMouseName(0)).toBe('MOUSE_LEFT')
    expect(engineMouseName(1)).toBe('MOUSE_MIDDLE')
    expect(engineMouseName(2)).toBe('MOUSE_RIGHT')
  })

  it('only produces names the engine knows', () => {
    for (let button = 0; button < 5; button++) {
      expect(ENGINE_KEYS.has(engineMouseName(button) as string)).toBe(true)
    }
    expect(engineMouseName(9)).toBeNull()
  })
})

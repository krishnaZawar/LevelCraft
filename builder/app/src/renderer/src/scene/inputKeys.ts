// Translates browser input events into the engine's key names.
//
// The engine's mapping.json carries GLFW key codes (Enter is 257, Left is
// 263), which do not line up with anything the browser reports, so the numeric
// code cannot be used to match them. The names do line up though, by a rule
// regular enough to derive rather than tabulate — and the backend's mapping
// decides which of those names actually exist (see runtimeStore).

// The codes whose engine name cannot be derived from their shape.
const NAMED_KEYS: Record<string, string> = {
  Space: 'KEY_SPACE',
  Escape: 'KEY_ESCAPE',
  Enter: 'KEY_ENTER',
  NumpadEnter: 'KEY_ENTER',
  Tab: 'KEY_TAB',
  Backspace: 'KEY_BACKSPACE',
  Insert: 'KEY_INSERT',
  Delete: 'KEY_DELETE',
  ArrowRight: 'KEY_RIGHT',
  ArrowLeft: 'KEY_LEFT',
  ArrowDown: 'KEY_DOWN',
  ArrowUp: 'KEY_UP',
  PageUp: 'KEY_PAGE_UP',
  PageDown: 'KEY_PAGE_DOWN',
  Home: 'KEY_HOME',
  End: 'KEY_END',
  ShiftLeft: 'KEY_LEFT_SHIFT',
  ControlLeft: 'KEY_LEFT_CONTROL',
  AltLeft: 'KEY_LEFT_ALT',
  MetaLeft: 'KEY_LEFT_SUPER',
  ShiftRight: 'KEY_RIGHT_SHIFT',
  ControlRight: 'KEY_RIGHT_CONTROL',
  AltRight: 'KEY_RIGHT_ALT',
  MetaRight: 'KEY_RIGHT_SUPER'
}

// MouseEvent.button, in the order the DOM numbers them.
const MOUSE_BUTTONS = ['MOUSE_LEFT', 'MOUSE_MIDDLE', 'MOUSE_RIGHT', 'MOUSE_BACK', 'MOUSE_FORWARD']

/**
 * The engine's name for a KeyboardEvent.code, or null if there is no such key.
 *
 * Returning a name is not a promise that the engine knows it — the caller
 * checks against the mapping it fetched.
 */
export function engineKeyName(code: string): string | null {
  if (/^Key[A-Z]$/.test(code)) return `KEY_${code.slice(3)}`
  if (/^Digit[0-9]$/.test(code)) return `KEY_${code.slice(5)}`
  if (/^F([1-9]|1[0-2])$/.test(code)) return `KEY_${code}`
  return NAMED_KEYS[code] ?? null
}

/** The engine's name for a MouseEvent.button, or null for a button it has none for. */
export function engineMouseName(button: number): string | null {
  return MOUSE_BUTTONS[button] ?? null
}

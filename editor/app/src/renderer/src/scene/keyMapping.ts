// No GET /input/schema endpoint exists on either backend yet (see
// docs/client/approach.md §4.3) — the symbolic names below are hand-copied
// from utils/input/mapping.json's actual `name` values for exactly the
// smoke-test subset needed to prove input round-trips through the WS
// pipeline. Deliberately not grown into a full keyboard map: once
// /input/schema exists, this should be replaced by fetching the real
// registry instead of extending this list further.
export const KEY_CODE_TO_NAME: Record<string, string> = {
  KeyW: 'KEY_W',
  KeyA: 'KEY_A',
  KeyS: 'KEY_S',
  KeyD: 'KEY_D',
  Space: 'KEY_SPACE',
  ArrowUp: 'KEY_UP',
  ArrowDown: 'KEY_DOWN',
  ArrowLeft: 'KEY_LEFT',
  ArrowRight: 'KEY_RIGHT'
}

import { useEffect, useRef, useState } from 'react'
import { Group, Layer, Rect, Stage } from 'react-konva'
import { AlertCircle, Loader2, Square } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { KEY_CODE_TO_NAME } from '@/scene/keyMapping'
import { computeFrameLayout, SCREEN_HEIGHT, SCREEN_WIDTH } from '@/scene/frameLayout'
import { renderGrid, renderScene } from '@/scene/renderScene'
import { useRuntimeStore } from '@/store/runtimeStore'

const FRAME_FILL = 'oklch(0.145 0 0)'
const FRAME_STROKE = 'oklch(1 0 0 / 10%)'
const GRID_CELL_SIZE = 64

function StatusBadge(): React.JSX.Element | null {
  const status = useRuntimeStore((state) => state.status)

  if (status === 'connecting') {
    return (
      <span className="text-muted-foreground flex items-center gap-1.5 text-xs">
        <Loader2 className="size-3.5 animate-spin" aria-hidden />
        Connecting…
      </span>
    )
  }
  if (status === 'connected') {
    return (
      <span className="flex items-center gap-1.5 text-xs">
        <span className="size-1.5 rounded-full bg-emerald-500" aria-hidden />
        <span className="text-muted-foreground">Connected</span>
      </span>
    )
  }
  if (status === 'error') {
    return (
      <span className="text-destructive flex items-center gap-1.5 text-xs">
        <AlertCircle className="size-3.5" aria-hidden />
        Connection error
      </span>
    )
  }
  return null
}

// Play mode: proves the connect → send input → receive snapshot → render
// pipeline (see docs/client/approach.md §4.1) — builder/backend has no
// gameplay wired up yet, so nothing will visibly react to input. Renders
// via the same renderScene() Workspace uses, fed the WS snapshot instead
// of editorStore, exactly why that function was kept data-source-agnostic.
function PlayView({ onStop }: { onStop: () => void }): React.JSX.Element {
  const status = useRuntimeStore((state) => state.status)
  const snapshot = useRuntimeStore((state) => state.snapshot)
  const error = useRuntimeStore((state) => state.error)
  const connect = useRuntimeStore((state) => state.connect)
  const sendKeyDown = useRuntimeStore((state) => state.sendKeyDown)
  const sendKeyUp = useRuntimeStore((state) => state.sendKeyUp)

  const containerRef = useRef<HTMLDivElement>(null)
  const [size, setSize] = useState({ width: 0, height: 0 })
  const heldKeys = useRef(new Set<string>())

  useEffect(() => {
    connect()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    const el = containerRef.current
    if (!el) return
    const observer = new ResizeObserver(([entry]) => {
      setSize({ width: entry.contentRect.width, height: entry.contentRect.height })
    })
    observer.observe(el)
    return () => observer.disconnect()
  }, [])

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent): void {
      // OS key-repeat fires keydown repeatedly while held — the transition-
      // based protocol wants exactly one KeyDownCommand per press (see
      // docs/builder/input_handling_and_modelling_design.md §9).
      if (e.repeat) return
      const keyName = KEY_CODE_TO_NAME[e.code]
      if (!keyName) return
      heldKeys.current.add(keyName)
      sendKeyDown(keyName)
    }

    function handleKeyUp(e: KeyboardEvent): void {
      const keyName = KEY_CODE_TO_NAME[e.code]
      if (!keyName) return
      heldKeys.current.delete(keyName)
      sendKeyUp(keyName)
    }

    // No INPUT_FOCUS_LOST/reset command exists server-side yet — if the
    // window loses focus while keys are held, the backend would keep
    // thinking they're down forever. Client-side workaround: flush every
    // key we've sent KeyDown for as soon as focus is lost.
    function handleBlur(): void {
      heldKeys.current.forEach((keyName) => sendKeyUp(keyName))
      heldKeys.current.clear()
    }

    window.addEventListener('keydown', handleKeyDown)
    window.addEventListener('keyup', handleKeyUp)
    window.addEventListener('blur', handleBlur)
    return () => {
      window.removeEventListener('keydown', handleKeyDown)
      window.removeEventListener('keyup', handleKeyUp)
      window.removeEventListener('blur', handleBlur)
      handleBlur()
    }
  }, [sendKeyDown, sendKeyUp])

  const hasSize = size.width > 0 && size.height > 0
  const { scale, offsetX, offsetY } = computeFrameLayout(size.width, size.height)

  return (
    <div className="bg-background text-foreground flex h-screen flex-col">
      <div className="border-border flex h-11 shrink-0 items-center justify-between border-b px-4">
        <span className="text-sm font-medium">Play Mode</span>
        <div className="flex items-center gap-3">
          <StatusBadge />
          <Button variant="destructive" size="sm" onClick={onStop}>
            <Square aria-hidden />
            Stop
          </Button>
        </div>
      </div>

      <div ref={containerRef} className="bg-workspace relative flex-1 overflow-hidden">
        {status === 'error' && (
          <div className="absolute inset-0 z-10 flex flex-col items-center justify-center gap-3 p-6 text-center">
            <div className="text-destructive flex items-center gap-2">
              <AlertCircle className="size-4" aria-hidden />
              <p className="text-sm font-medium">Builder unreachable</p>
            </div>
            {error && <p className="text-muted-foreground max-w-xs text-xs">{error}</p>}
            <Button variant="secondary" size="sm" onClick={connect}>
              Retry
            </Button>
          </div>
        )}

        {hasSize && (
          <Stage width={size.width} height={size.height}>
            <Layer>
              <Group x={offsetX} y={offsetY} scaleX={scale} scaleY={scale}>
                <Rect
                  x={0}
                  y={0}
                  width={SCREEN_WIDTH}
                  height={SCREEN_HEIGHT}
                  fill={FRAME_FILL}
                  stroke={FRAME_STROKE}
                  strokeWidth={1 / scale}
                />
                {renderGrid(SCREEN_WIDTH, SCREEN_HEIGHT, GRID_CELL_SIZE, scale)}
                {snapshot && renderScene(snapshot)}
              </Group>
            </Layer>
          </Stage>
        )}

        {!snapshot && status === 'connected' && (
          <p className="text-muted-foreground pointer-events-none absolute inset-0 flex items-center justify-center text-xs">
            Waiting for the first scene snapshot (every ~5s)…
          </p>
        )}
      </div>
    </div>
  )
}

export default PlayView

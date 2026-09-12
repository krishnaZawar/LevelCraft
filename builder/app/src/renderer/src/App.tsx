import { useEffect, useRef, useState } from 'react'
import { Group, Layer, Rect, Stage } from 'react-konva'
import { computeFrameLayout, SCREEN_HEIGHT, SCREEN_WIDTH } from '@/scene/frameLayout'
import { KEY_CODE_TO_NAME } from '@/scene/keyMapping'
import { renderScene } from '@/scene/renderScene'
import { useRuntimeStore } from '@/store/runtimeStore'

const FRAME_FILL = '#121212'

function StatusBadge(): React.JSX.Element {
  const status = useRuntimeStore((state) => state.status)

  const label =
    status === 'connected' ? 'Connected' : status === 'connecting' ? 'Connecting…' : 'Disconnected'

  return (
    <span className="player__status">
      <span className={`player__dot player__dot--${status}`} aria-hidden />
      {label}
    </span>
  )
}

// The player: connect, forward input transitions, render what comes back.
// No gameplay is wired up backend-side yet, so this proves the pipeline only.
function App(): React.JSX.Element {
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
      // OS key-repeat fires keydown repeatedly while held, but the protocol
      // wants exactly one KeyDownCommand per press.
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

    // Without a server-side INPUT_FOCUS_LOST, keys held at blur would stay
    // down forever — so flush every key we have sent a KeyDown for.
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
    <div className="player">
      <div className="player__bar">
        <span className="player__title">LevelCraft</span>
        <StatusBadge />
      </div>

      <div ref={containerRef} className="player__stage">
        {status === 'error' && (
          <div className="player__overlay">
            <p className="player__overlay-title">Disconnected from the game</p>
            {error && <p className="player__overlay-body">{error}</p>}
            <button type="button" className="player__button" onClick={connect}>
              Retry
            </button>
          </div>
        )}

        {hasSize && (
          <Stage width={size.width} height={size.height}>
            <Layer>
              <Group x={offsetX} y={offsetY} scaleX={scale} scaleY={scale}>
                <Rect x={0} y={0} width={SCREEN_WIDTH} height={SCREEN_HEIGHT} fill={FRAME_FILL} />
                {snapshot && renderScene(snapshot)}
              </Group>
            </Layer>
          </Stage>
        )}
      </div>
    </div>
  )
}

export default App

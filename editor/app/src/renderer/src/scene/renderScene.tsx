import { Line, Rect } from 'react-konva'
import type Konva from 'konva'
import type { GameState } from '@/api/gameApi'

// Deliberately independent of editorStore: Phase 7 (Play mode) feeds this
// the same shape from the builder's WebSocket snapshots instead, and
// shouldn't need this logic rewritten to do so.
interface RenderSceneOptions {
  selectedObjectId?: string | null
  onSelect?: (id: string) => void
  onDragEnd?: (id: string, x: number, y: number) => void
  // Fired when a Transformer handle (see Workspace, which owns the actual
  // Transformer node) finishes a resize on this shape.
  onTransformEnd?: (id: string, node: Konva.Node) => void
  draggable?: boolean
  // The scale of the enclosing Group (e.g. a zoom-to-fit transform) — divided
  // into strokeWidth so it stays a constant physical size regardless of how
  // zoomed in/out the scene is rendered.
  scale?: number
}

export const SELECTED_STROKE = 'oklch(0.74 0.14 55)'
const DEFAULT_FILL = 'oklch(0.269 0 0)'
// A faint outline on every object, selected or not — flat, borderless color
// blocks read as placeholder/toy shapes; real engines always give sprites
// an edge against the canvas.
const OBJECT_STROKE = 'oklch(0 0 0 / 35%)'

function colorToCss(details: Record<string, unknown> | undefined): string {
  if (!details) return DEFAULT_FILL
  const r = Number(details.r ?? 0)
  const g = Number(details.g ?? 0)
  const b = Number(details.b ?? 0)
  const a = Number(details.a ?? 255) / 255
  return `rgba(${r}, ${g}, ${b}, ${a})`
}

// Each Rect gets a Konva id matching the gameobject's id, so Workspace can
// find the selected node (stage.findOne('#' + id)) to attach its Transformer
// to — the actual drag-to-resize handles live there, not in this function.
export function renderScene(
  gameObjects: GameState,
  options: RenderSceneOptions = {}
): React.ReactNode[] {
  const {
    selectedObjectId,
    onSelect,
    onDragEnd,
    onTransformEnd,
    draggable = false,
    scale = 1
  } = options

  return Object.values(gameObjects)
    .filter((object) => object.components.Transform)
    .map((object) => {
      const transform = object.components.Transform
      const x = Number(transform.x ?? 0)
      const y = Number(transform.y ?? 0)
      const width = Number(transform.w ?? 0)
      const height = Number(transform.h ?? 0)
      const isSelected = object.id === selectedObjectId

      return (
        <Rect
          key={object.id}
          id={object.id}
          x={x}
          y={y}
          width={width}
          height={height}
          fill={colorToCss(object.components.Color)}
          stroke={isSelected ? SELECTED_STROKE : OBJECT_STROKE}
          strokeWidth={(isSelected ? 2 : 1) / scale}
          draggable={draggable}
          onClick={() => onSelect?.(object.id)}
          onTap={() => onSelect?.(object.id)}
          onDragEnd={(e: Konva.KonvaEventObject<DragEvent>) =>
            onDragEnd?.(object.id, e.target.x(), e.target.y())
          }
          onTransformEnd={(e: Konva.KonvaEventObject<Event>) =>
            onTransformEnd?.(object.id, e.target)
          }
        />
      )
    })
}

// A faint reference grid inside the screen frame, spaced in world units —
// the kind of alignment grid every 2D scene view (Unity, Godot, Figma
// frames) shows, as opposed to the pasteboard's own (removed) dot pattern.
export function renderGrid(
  width: number,
  height: number,
  cellSize: number,
  scale: number
): React.ReactNode[] {
  const lines: React.ReactNode[] = []
  const strokeWidth = 1 / scale

  for (let x = 0; x <= width; x += cellSize) {
    lines.push(
      <Line
        key={`v-${x}`}
        points={[x, 0, x, height]}
        stroke="oklch(1 0 0 / 5%)"
        strokeWidth={strokeWidth}
        listening={false}
      />
    )
  }
  for (let y = 0; y <= height; y += cellSize) {
    lines.push(
      <Line
        key={`h-${y}`}
        points={[0, y, width, y]}
        stroke="oklch(1 0 0 / 5%)"
        strokeWidth={strokeWidth}
        listening={false}
      />
    )
  }
  return lines
}

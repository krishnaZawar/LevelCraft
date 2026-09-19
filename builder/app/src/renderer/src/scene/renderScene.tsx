import { Rect } from 'react-konva'
import type { GameState } from './types'

const DEFAULT_FILL = 'rgb(69, 69, 69)'

function colorToCss(details: Record<string, unknown> | undefined): string {
  if (!details) return DEFAULT_FILL
  const r = Number(details.r ?? 0)
  const g = Number(details.g ?? 0)
  const b = Number(details.b ?? 0)
  const a = Number(details.a ?? 255) / 255
  return `rgba(${r}, ${g}, ${b}, ${a})`
}

// Draws the scene as the player sees it: no selection outlines, no drag
// handles, no editor affordances — only what the game state is.
export function renderScene(gameObjects: GameState): React.ReactNode[] {
  return Object.values(gameObjects)
    .filter((object) => object.components.Transform)
    .map((object) => {
      const transform = object.components.Transform
      return (
        <Rect
          key={object.id}
          x={Number(transform.x ?? 0)}
          y={Number(transform.y ?? 0)}
          width={Number(transform.w ?? 0)}
          height={Number(transform.h ?? 0)}
          fill={colorToCss(object.components.Color)}
          listening={false}
        />
      )
    })
}

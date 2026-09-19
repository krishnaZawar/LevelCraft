// Frames the game screen the same way the editor's Workspace panel does, so a
// scene lays out identically whether it is being authored or played.

// Default game screen size, matching a standard 16:9 HD canvas.
export const SCREEN_WIDTH = 1280
export const SCREEN_HEIGHT = 720
export const SCREEN_PADDING = 24

export interface FrameLayout {
  scale: number
  offsetX: number
  offsetY: number
  frameWidth: number
  frameHeight: number
}

// Zoom-to-fits the screen into the container and centers it, never scaling
// past 1:1 so a large window doesn't blow the frame up unrealistically.
export function computeFrameLayout(containerWidth: number, containerHeight: number): FrameLayout {
  const hasSize = containerWidth > 0 && containerHeight > 0
  const scale = hasSize
    ? Math.min(
        (containerWidth - SCREEN_PADDING * 2) / SCREEN_WIDTH,
        (containerHeight - SCREEN_PADDING * 2) / SCREEN_HEIGHT,
        1
      )
    : 1
  const frameWidth = SCREEN_WIDTH * scale
  const frameHeight = SCREEN_HEIGHT * scale
  return {
    scale,
    frameWidth,
    frameHeight,
    offsetX: (containerWidth - frameWidth) / 2,
    offsetY: (containerHeight - frameHeight) / 2
  }
}

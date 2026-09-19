// Frames the "Main Screen" artboard the Workspace panel authors against.
// builder/app keeps its own matching copy: the two apps agree, not share.

// Default game screen size, matching a standard 16:9 HD canvas. No pan/zoom
// yet (not called for by the current phase) — the frame just zoom-to-fits
// the panel each time it resizes.
export const SCREEN_WIDTH = 1280
export const SCREEN_HEIGHT = 720
export const SCREEN_PADDING = 48

export interface FrameLayout {
  scale: number
  offsetX: number
  offsetY: number
  frameWidth: number
  frameHeight: number
}

// Zoom-to-fits SCREEN_WIDTH x SCREEN_HEIGHT into the given container size,
// never scaling up past 1:1 so a large panel doesn't blow the frame up
// unrealistically, and centers it within the container.
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

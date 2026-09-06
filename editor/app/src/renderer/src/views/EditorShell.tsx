import { useEffect, useRef, useState } from 'react'
import { Group, Layer, Rect, Stage, Transformer } from 'react-konva'
import type Konva from 'konva'
import AppMenu from '@/AppMenu'
import {
  AlertCircle,
  Box,
  Copy,
  LayoutGrid,
  ListTree,
  Loader2,
  Play,
  Plus,
  Save,
  SlidersHorizontal,
  Trash2,
  X
} from 'lucide-react'
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { toast } from 'sonner'
import { saveGame } from '@/api/gameApi'
import { cn } from '@/lib/utils'
import { computeFrameLayout, SCREEN_HEIGHT, SCREEN_WIDTH } from '@/scene/frameLayout'
import { renderGrid, renderScene, SELECTED_STROKE } from '@/scene/renderScene'
import { useEditorStore } from '@/store/editorStore'
import { useProjectStore } from '@/store/projectStore'

// Keeps field order stable and predictable instead of relying on JSON key
// order. Falls back to whatever keys the component actually has, so a
// future component type still renders without needing an entry here.
const COMPONENT_FIELD_ORDER: Record<string, string[]> = {
  Transform: ['x', 'y', 'w', 'h'],
  Color: ['r', 'g', 'b', 'a']
}

// Matches utils/component/base/const.go's ColorValueRangeMin/Max — clamped
// client-side so an out-of-range value never reaches the backend at all,
// rather than round-tripping to find out it gets rejected.
const COMPONENT_FIELD_RANGE: Record<string, Record<string, [number, number]>> = {
  Color: { r: [0, 255], g: [0, 255], b: [0, 255], a: [0, 255] }
}

const FRAME_FILL = 'oklch(0.145 0 0)'
const FRAME_STROKE = 'oklch(1 0 0 / 10%)'
const GRID_CELL_SIZE = 64

function PanelHeading({
  icon: Icon,
  children
}: {
  icon: React.ComponentType<{ className?: string }>
  children: React.ReactNode
}): React.JSX.Element {
  return (
    <h2 className="text-muted-foreground border-border flex items-center gap-1.5 border-b px-3 py-2 text-xs font-semibold tracking-wide uppercase">
      <Icon className="size-3.5" />
      {children}
    </h2>
  )
}

function TopBar(): React.JSX.Element {
  const activeProject = useProjectStore((state) => state.activeProject)
  const closeProject = useProjectStore((state) => state.closeProject)

  return (
    <div
      className={cn(
        'border-border flex h-11 shrink-0 items-center justify-between border-b pr-3',
        window.api.platform === 'darwin' ? 'pl-20' : 'pl-3'
      )}
      style={{ WebkitAppRegion: 'drag' } as React.CSSProperties}
    >
      <div
        className="flex items-center gap-2"
        style={{ WebkitAppRegion: 'no-drag' } as React.CSSProperties}
      >
        <AppMenu />
        <span className="text-sm font-medium">{activeProject?.name}</span>
      </div>
      <Button
        variant="ghost"
        size="sm"
        onClick={closeProject}
        style={{ WebkitAppRegion: 'no-drag' } as React.CSSProperties}
      >
        <X aria-hidden />
        Close Project
      </Button>
    </div>
  )
}

function HierarchyPanel(): React.JSX.Element {
  const gameObjects = useEditorStore((state) => state.gameObjects)
  const objectOrder = useEditorStore((state) => state.objectOrder)
  const isLoading = useEditorStore((state) => state.isLoading)
  const loadError = useEditorStore((state) => state.loadError)
  const selectedObjectId = useEditorStore((state) => state.selectedObjectId)
  const selectObject = useEditorStore((state) => state.selectObject)
  const reorderObjects = useEditorStore((state) => state.reorderObjects)
  const fetchGameState = useEditorStore((state) => state.fetchGameState)

  const [draggedId, setDraggedId] = useState<string | null>(null)

  const objects = objectOrder.map((id) => gameObjects[id]).filter(Boolean)

  if (loadError) {
    return (
      <div className="flex flex-1 flex-col items-start gap-2 p-3">
        <div className="text-destructive flex items-center gap-2">
          <AlertCircle className="size-4" aria-hidden />
          <p className="text-sm font-medium">Couldn&apos;t load the scene</p>
        </div>
        <Button variant="secondary" size="sm" onClick={fetchGameState}>
          Try again
        </Button>
      </div>
    )
  }

  if (isLoading && objects.length === 0) {
    return (
      <div className="text-muted-foreground flex flex-1 flex-col items-center justify-center gap-2 p-6 text-center text-xs">
        <Loader2 className="size-4 animate-spin" aria-hidden />
        <p>Loading scene…</p>
      </div>
    )
  }

  if (objects.length === 0) {
    return (
      <div className="text-muted-foreground flex flex-1 flex-col items-center justify-center gap-1 p-6 text-center text-xs">
        <Box className="size-5" aria-hidden />
        <p>No objects yet</p>
      </div>
    )
  }

  return (
    <ScrollArea className="flex-1">
      <ul className="p-1">
        {objects.map((object) => {
          const isSelected = selectedObjectId === object.id
          const color = object.components.Color

          return (
            <li key={object.id}>
              <button
                type="button"
                draggable
                onClick={() => selectObject(object.id)}
                onDragStart={(e) => {
                  // Read the dragged id back out of dataTransfer in onDrop
                  // rather than this component's own state — state set here
                  // isn't guaranteed to have flushed by the time onDrop fires.
                  e.dataTransfer.setData('text/plain', object.id)
                  setDraggedId(object.id)
                }}
                onDragEnd={() => setDraggedId(null)}
                onDragOver={(e) => e.preventDefault()}
                onDrop={(e) => {
                  e.preventDefault()
                  const sourceId = e.dataTransfer.getData('text/plain')
                  if (sourceId) reorderObjects(sourceId, object.id)
                  setDraggedId(null)
                }}
                className={cn(
                  'flex w-full items-center gap-2 rounded-sm border-l-2 py-1.5 pr-2 pl-[calc(0.5rem-2px)] text-left text-sm transition-colors duration-100',
                  draggedId === object.id && 'opacity-40',
                  isSelected
                    ? 'bg-sidebar-primary/15 border-sidebar-primary text-foreground'
                    : 'hover:bg-accent hover:text-accent-foreground border-transparent'
                )}
              >
                <Box
                  className={cn('size-3.5 shrink-0', isSelected && 'text-sidebar-primary')}
                  aria-hidden
                />
                <span className="truncate">{object.name || 'GameObject'}</span>
                {color && (
                  <span
                    className="border-border ml-auto size-2.5 shrink-0 rounded-[2px] border"
                    style={{
                      backgroundColor: `rgba(${Number(color.r ?? 0)}, ${Number(color.g ?? 0)}, ${Number(color.b ?? 0)}, ${Number(color.a ?? 255) / 255})`
                    }}
                    aria-hidden
                  />
                )}
              </button>
            </li>
          )
        })}
      </ul>
    </ScrollArea>
  )
}

function NumberField({
  label,
  value,
  onCommit,
  range
}: {
  label: string
  value: number
  onCommit: (value: number) => void
  range?: [number, number]
}): React.JSX.Element {
  // Re-initializes from `value` whenever it changes externally: the parent
  // remounts this via a `value`-derived key instead of syncing with an effect.
  const [text, setText] = useState(String(value))
  const isInvalid = text.trim() !== '' && !Number.isFinite(Number(text))

  function commit(): void {
    let parsed = Number(text)
    if (!Number.isFinite(parsed)) {
      setText(String(value))
      return
    }
    // Clamp rather than reject outright — matches how Unity/most inspectors
    // handle out-of-range numeric input instead of silently sending it
    // through and finding out the backend rejected it.
    if (range) {
      parsed = Math.min(range[1], Math.max(range[0], Math.round(parsed)))
    }
    if (parsed !== value) {
      onCommit(parsed)
    }
    setText(String(parsed))
  }

  return (
    <label className="flex items-center gap-1.5 text-xs">
      <span className="text-muted-foreground w-3 shrink-0 uppercase">{label}</span>
      <Input
        value={text}
        onChange={(e) => setText(e.target.value)}
        onBlur={commit}
        onKeyDown={(e) => {
          if (e.key === 'Enter') (e.target as HTMLInputElement).blur()
        }}
        inputMode="numeric"
        aria-invalid={isInvalid}
        className="h-7 font-mono text-xs tabular-nums"
      />
    </label>
  )
}

function ComponentCard({
  objectId,
  name,
  details
}: {
  objectId: string
  name: string
  details: Record<string, unknown>
}): React.JSX.Element {
  const updateComponent = useEditorStore((state) => state.updateComponent)
  const deleteComponent = useEditorStore((state) => state.deleteComponent)

  const fields = COMPONENT_FIELD_ORDER[name] ?? Object.keys(details)

  return (
    <div className="border-border space-y-2 rounded-md border p-2.5">
      <div className="flex items-center justify-between">
        <span className="text-xs font-semibold">{name}</span>
        <Button
          variant="ghost"
          size="icon"
          className="size-6"
          onClick={() => deleteComponent(objectId, name)}
        >
          <Trash2 className="size-3.5" aria-hidden />
        </Button>
      </div>
      <div className="grid grid-cols-2 gap-x-3 gap-y-1.5">
        {fields.map((field) => {
          const value = Number(details[field] ?? 0)
          return (
            <NumberField
              key={`${objectId}:${field}:${value}`}
              label={field}
              value={value}
              range={COMPONENT_FIELD_RANGE[name]?.[field]}
              onCommit={(next) => updateComponent(objectId, name, { ...details, [field]: next })}
            />
          )
        })}
      </div>
    </div>
  )
}

function AttributesPanel(): React.JSX.Element {
  const gameObjects = useEditorStore((state) => state.gameObjects)
  const selectedObjectId = useEditorStore((state) => state.selectedObjectId)
  const availableComponents = useEditorStore((state) => state.availableComponents)
  const fetchAvailableComponents = useEditorStore((state) => state.fetchAvailableComponents)
  const addComponent = useEditorStore((state) => state.addComponent)

  useEffect(() => {
    fetchAvailableComponents()
  }, [fetchAvailableComponents])

  const selected = selectedObjectId ? gameObjects[selectedObjectId] : null

  if (!selected) {
    return (
      <div className="text-muted-foreground flex flex-1 flex-col items-center justify-center gap-1 p-6 text-center text-xs">
        <SlidersHorizontal className="size-5" aria-hidden />
        <p>No object selected</p>
      </div>
    )
  }

  const attachedNames = Object.keys(selected.components)
  const addableComponents = availableComponents.filter((name) => !attachedNames.includes(name))

  return (
    <div className="flex flex-1 flex-col overflow-hidden">
      <div className="border-border space-y-0.5 border-b p-3">
        <p className="truncate text-sm font-medium">{selected.name || 'GameObject'}</p>
        {selected.group && <p className="text-muted-foreground text-xs">{selected.group}</p>}
      </div>

      <ScrollArea className="flex-1">
        <div className="space-y-2 p-3">
          {attachedNames.length === 0 ? (
            <p className="text-muted-foreground text-xs">No components yet.</p>
          ) : (
            attachedNames.map((name) => (
              <ComponentCard
                key={name}
                objectId={selected.id}
                name={name}
                details={selected.components[name]}
              />
            ))
          )}
        </div>
      </ScrollArea>

      {addableComponents.length > 0 && (
        <div className="border-border border-t p-2">
          <Select value="" onValueChange={(name) => addComponent(selected.id, name)}>
            <SelectTrigger size="sm" className="w-full">
              <SelectValue placeholder="Add component..." />
            </SelectTrigger>
            <SelectContent>
              {addableComponents.map((name) => (
                <SelectItem key={name} value={name}>
                  {name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      )}
    </div>
  )
}

function UtilityPanel(): React.JSX.Element {
  const selectedObjectId = useEditorStore((state) => state.selectedObjectId)
  const addGameobject = useEditorStore((state) => state.addGameobject)
  const deleteGameobject = useEditorStore((state) => state.deleteGameobject)
  const activeProject = useProjectStore((state) => state.activeProject)

  const [saveState, setSaveState] = useState<'idle' | 'saving' | 'saved'>('idle')
  // Play mode runs in its own pop-out window/process (see main/index.ts's
  // createPlayWindow), not inline here — this just tracks whether one is
  // open, so Run stays disabled for its whole duration (the backend port
  // is fixed, so a second launch would collide with it).
  const [isRunning, setIsRunning] = useState(false)

  useEffect(() => {
    window.api.builder.onStopped(() => setIsRunning(false))
  }, [])

  async function handleSave(): Promise<void> {
    if (!activeProject) return
    setSaveState('saving')
    try {
      await saveGame(activeProject.scenePath)
      setSaveState('saved')
      setTimeout(() => setSaveState('idle'), 1500)
    } catch (err) {
      setSaveState('idle')
      toast.error("Couldn't save the scene", { description: (err as Error).message })
    }
  }

  // Per docs/client/phases.md Phase 7: save the current scene to the
  // project's real scene file (no separate scratch path needed now that
  // one exists), then hand that same path to builder/backend.
  async function handleRun(): Promise<void> {
    if (!activeProject || isRunning) return
    setIsRunning(true)
    try {
      await saveGame(activeProject.scenePath)
    } catch (err) {
      setIsRunning(false)
      toast.error("Couldn't save the scene", { description: (err as Error).message })
      return
    }
    const result = await window.api.builder.launch(activeProject.scenePath)
    if (!result.ok) {
      setIsRunning(false)
      toast.error("Couldn't start builder/backend", { description: result.message })
    }
  }

  return (
    <div className="flex flex-1 flex-col gap-3 p-3">
      <div className="flex flex-wrap gap-2">
        <Button variant="secondary" size="sm" onClick={addGameobject}>
          <Plus aria-hidden />
          Add Object
        </Button>

        <Button
          variant="outline"
          size="sm"
          disabled={!selectedObjectId}
          onClick={() => selectedObjectId && deleteGameobject(selectedObjectId)}
        >
          <Trash2 aria-hidden />
          Delete
        </Button>

        <Tooltip>
          <TooltipTrigger asChild>
            <span>
              <Button variant="outline" size="sm" disabled>
                <Copy aria-hidden />
                Duplicate
              </Button>
            </span>
          </TooltipTrigger>
          <TooltipContent>Pending implementation</TooltipContent>
        </Tooltip>

        <Button
          variant="outline"
          size="sm"
          onClick={handleSave}
          disabled={!activeProject || saveState === 'saving'}
        >
          <Save aria-hidden />
          {saveState === 'saving' ? 'Saving…' : saveState === 'saved' ? 'Saved' : 'Save'}
        </Button>
      </div>

      <Button className="w-full" onClick={handleRun} disabled={!activeProject || isRunning}>
        <Play aria-hidden />
        {isRunning ? 'Running…' : 'Run'}
      </Button>
    </div>
  )
}

function Workspace(): React.JSX.Element {
  const gameObjects = useEditorStore((state) => state.gameObjects)
  const selectedObjectId = useEditorStore((state) => state.selectedObjectId)
  const selectObject = useEditorStore((state) => state.selectObject)
  const updateComponent = useEditorStore((state) => state.updateComponent)

  const containerRef = useRef<HTMLDivElement>(null)
  const stageRef = useRef<Konva.Stage>(null)
  const transformerRef = useRef<Konva.Transformer>(null)
  const [size, setSize] = useState({ width: 0, height: 0 })

  useEffect(() => {
    const el = containerRef.current
    if (!el) return

    const observer = new ResizeObserver(([entry]) => {
      setSize({ width: entry.contentRect.width, height: entry.contentRect.height })
    })
    observer.observe(el)
    return () => observer.disconnect()
  }, [])

  // Attaches the Transformer's drag-to-resize handles to whichever shape is
  // selected — re-runs on gameObjects changes too, since the node for a
  // freshly-selected object may not exist yet if it has no Transform.
  useEffect(() => {
    const transformer = transformerRef.current
    const stage = stageRef.current
    if (!transformer || !stage) return

    const node = selectedObjectId ? stage.findOne(`#${selectedObjectId}`) : null
    transformer.nodes(node ? [node] : [])
    transformer.getLayer()?.batchDraw()
  }, [selectedObjectId, gameObjects])

  function handleDragEnd(objectId: string, x: number, y: number): void {
    const transform = gameObjects[objectId]?.components.Transform
    if (!transform) return
    updateComponent(objectId, 'Transform', { ...transform, x: Math.round(x), y: Math.round(y) })
  }

  function handleTransformEnd(objectId: string, node: Konva.Node): void {
    const transform = gameObjects[objectId]?.components.Transform
    if (!transform) return

    // Transformer resizes via scale, not width/height directly — bake the
    // scale into width/height and reset it, so Transform.w/h stay the
    // source of truth instead of drifting out of sync with a lingering scale.
    const x = Math.round(node.x())
    const y = Math.round(node.y())
    const width = Math.max(1, Math.round(node.width() * node.scaleX()))
    const height = Math.max(1, Math.round(node.height() * node.scaleY()))

    // Applying these as separate node.scaleX(1)/scaleY(1) calls (before
    // width/height ever change) makes Konva redraw the node at its old
    // size for one frame first — a visible flash back to the pre-resize
    // shape. setAttrs applies all of them atomically in the same redraw.
    node.setAttrs({ x, y, width, height, scaleX: 1, scaleY: 1 })

    updateComponent(objectId, 'Transform', { ...transform, x, y, w: width, h: height })
  }

  function handleStageMouseDown(e: Konva.KonvaEventObject<MouseEvent>): void {
    // Clicked empty pasteboard, not a shape — clear selection.
    if (e.target === e.target.getStage()) selectObject(null)
  }

  const hasSize = size.width > 0 && size.height > 0
  const { scale, offsetX, offsetY } = computeFrameLayout(size.width, size.height)

  return (
    <div ref={containerRef} className="bg-workspace relative h-full overflow-hidden">
      <h2 className="sr-only">Workspace</h2>
      {hasSize && (
        <>
          <span
            className="text-muted-foreground pointer-events-none absolute z-10 text-xs"
            style={{ left: offsetX, top: offsetY - 22 }}
          >
            Main Screen · {SCREEN_WIDTH}×{SCREEN_HEIGHT}
          </span>
          <Stage
            ref={stageRef}
            width={size.width}
            height={size.height}
            onMouseDown={handleStageMouseDown}
          >
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
                {renderScene(gameObjects, {
                  selectedObjectId,
                  onSelect: selectObject,
                  onDragEnd: handleDragEnd,
                  onTransformEnd: handleTransformEnd,
                  draggable: true,
                  scale
                })}
                <Transformer
                  ref={transformerRef}
                  rotateEnabled={false}
                  anchorSize={8 / scale}
                  anchorCornerRadius={2 / scale}
                  borderStroke={SELECTED_STROKE}
                  anchorStroke={SELECTED_STROKE}
                  anchorFill={FRAME_FILL}
                  borderStrokeWidth={1 / scale}
                  anchorStrokeWidth={1.5 / scale}
                  boundBoxFunc={(oldBox, newBox) =>
                    newBox.width < 10 || newBox.height < 10 ? oldBox : newBox
                  }
                />
              </Group>
            </Layer>
          </Stage>
        </>
      )}
    </div>
  )
}

function EditorShell(): React.JSX.Element {
  const fetchGameState = useEditorStore((state) => state.fetchGameState)

  useEffect(() => {
    fetchGameState()
  }, [fetchGameState])

  return (
    <div className="bg-background text-foreground flex h-screen flex-col">
      <TopBar />
      <div className="min-h-0 flex-1">
        <ResizablePanelGroup orientation="horizontal">
          <ResizablePanel defaultSize="18%" minSize="12%" maxSize="30%">
            <div className="flex h-full flex-col overflow-hidden">
              <PanelHeading icon={ListTree}>Hierarchy</PanelHeading>
              <HierarchyPanel />
            </div>
          </ResizablePanel>

          <ResizableHandle withHandle />

          <ResizablePanel defaultSize="62%">
            <ResizablePanelGroup orientation="vertical">
              <ResizablePanel defaultSize="78%" minSize="40%">
                <Workspace />
              </ResizablePanel>

              <ResizableHandle withHandle />

              <ResizablePanel defaultSize="22%" minSize="12%" maxSize="45%">
                <div className="flex h-full flex-col overflow-hidden">
                  <PanelHeading icon={LayoutGrid}>Utility</PanelHeading>
                  <UtilityPanel />
                </div>
              </ResizablePanel>
            </ResizablePanelGroup>
          </ResizablePanel>

          <ResizableHandle withHandle />

          <ResizablePanel defaultSize="20%" minSize="15%" maxSize="35%">
            <div className="flex h-full flex-col overflow-hidden">
              <PanelHeading icon={SlidersHorizontal}>Attributes</PanelHeading>
              <AttributesPanel />
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  )
}

export default EditorShell

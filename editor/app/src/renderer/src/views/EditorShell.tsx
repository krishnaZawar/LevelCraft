import { useCallback, useEffect, useRef, useState } from 'react'
import { Layer, Rect, Stage, Transformer } from 'react-konva'
import type Konva from 'konva'
import AppMenu from '@/AppMenu'
import {
  AlertCircle,
  Box,
  ChevronRight,
  Copy,
  GripVertical,
  ListTree,
  Loader2,
  Maximize2,
  Play,
  Plus,
  Save,
  SlidersHorizontal,
  Square,
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
import { isOrchestratorAvailable, ORCHESTRATOR_UNAVAILABLE_MESSAGE } from '@/api/orchestratorApi'
import { cn } from '@/lib/utils'
import { computeFrameLayout, SCREEN_HEIGHT, SCREEN_WIDTH } from '@/scene/frameLayout'
import { renderGrid, renderScene, SELECTED_STROKE } from '@/scene/renderScene'
import { componentNamesInOrder, useEditorStore } from '@/store/editorStore'
import { useProjectStore } from '@/store/projectStore'
import { RunStatus, useRunStore } from '@/store/runStore'
import { ProjectSummary } from '../../../shared/project'

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
// A solid white edge, so the game screen reads as a definite boundary on an
// otherwise unbounded canvas rather than as one more faint panel edge.
const SCREEN_STROKE = 'oklch(1 0 0 / 85%)'
const SCREEN_NODE_NAME = 'game-screen'
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

function TextField({
  label,
  value,
  placeholder,
  onCommit
}: {
  label: string
  value: string
  placeholder: string
  onCommit: (value: string) => void
}): React.JSX.Element {
  // Same remount-on-external-change approach NumberField uses, rather than
  // syncing with an effect.
  const [text, setText] = useState(value)

  function commit(): void {
    const next = text.trim()
    if (next === value) return
    onCommit(next)
  }

  return (
    <label className="flex flex-col gap-1">
      <span className="text-muted-foreground text-[10px] tracking-wide uppercase">{label}</span>
      <Input
        value={text}
        placeholder={placeholder}
        onChange={(e) => setText(e.target.value)}
        onBlur={commit}
        onKeyDown={(e) => {
          if (e.key === 'Enter') e.currentTarget.blur()
          // Abandoning the edit restores what was there, so a half-typed name
          // is never committed by clicking away afterwards.
          if (e.key === 'Escape') {
            setText(value)
            e.currentTarget.blur()
          }
        }}
        className="h-7 text-xs"
      />
    </label>
  )
}

function ComponentCard({
  objectId,
  name,
  details,
  isDragging,
  onDragStart,
  onDragEnd,
  onDrop
}: {
  objectId: string
  name: string
  details: Record<string, unknown>
  isDragging: boolean
  onDragStart: (name: string) => void
  onDragEnd: () => void
  onDrop: (sourceName: string, targetName: string) => void
}): React.JSX.Element {
  const updateComponent = useEditorStore((state) => state.updateComponent)
  const deleteComponent = useEditorStore((state) => state.deleteComponent)
  const isCollapsed = useEditorStore((state) => state.collapsedComponents.includes(name))
  const toggleCollapsed = useEditorStore((state) => state.toggleComponentCollapsed)

  const fields = COMPONENT_FIELD_ORDER[name] ?? Object.keys(details)

  return (
    <div
      draggable
      onDragStart={(e) => {
        e.dataTransfer.setData('text/plain', name)
        onDragStart(name)
      }}
      onDragEnd={onDragEnd}
      onDragOver={(e) => e.preventDefault()}
      onDrop={(e) => {
        e.preventDefault()
        // The dragged name comes back out of dataTransfer rather than from the
        // isDragging state, which isn't guaranteed to have flushed by the time
        // this fires — the same reason the Hierarchy reads it this way.
        const sourceName = e.dataTransfer.getData('text/plain')
        if (sourceName) onDrop(sourceName, name)
        onDragEnd()
      }}
      className={cn(
        'border-border rounded-md border transition-opacity',
        isDragging && 'opacity-40'
      )}
    >
      <div className="flex items-center gap-1 p-1.5 pl-2">
        <GripVertical className="text-muted-foreground size-3.5 shrink-0 cursor-grab" aria-hidden />
        <button
          type="button"
          onClick={() => toggleCollapsed(name)}
          aria-expanded={!isCollapsed}
          className="flex flex-1 items-center gap-1 text-left text-xs font-semibold"
        >
          <ChevronRight
            className={cn('size-3.5 transition-transform', !isCollapsed && 'rotate-90')}
            aria-hidden
          />
          {name}
        </button>
        <Button
          variant="ghost"
          size="icon"
          className="size-6"
          onClick={() => deleteComponent(objectId, name)}
        >
          <Trash2 className="size-3.5" aria-hidden />
        </Button>
      </div>

      {!isCollapsed && (
        <div className="grid grid-cols-2 gap-x-3 gap-y-1.5 px-2.5 pt-0.5 pb-2.5">
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
      )}
    </div>
  )
}

function AttributesPanel(): React.JSX.Element {
  const gameObjects = useEditorStore((state) => state.gameObjects)
  const selectedObjectId = useEditorStore((state) => state.selectedObjectId)
  const availableComponents = useEditorStore((state) => state.availableComponents)
  const fetchAvailableComponents = useEditorStore((state) => state.fetchAvailableComponents)
  const addComponent = useEditorStore((state) => state.addComponent)
  const componentOrder = useEditorStore((state) => state.componentOrder)
  const reorderComponents = useEditorStore((state) => state.reorderComponents)
  const updateMetadata = useEditorStore((state) => state.updateGameobjectMetadata)

  const [draggedComponent, setDraggedComponent] = useState<string | null>(null)

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

  const attachedNames = componentNamesInOrder(
    componentOrder[selected.id],
    Object.keys(selected.components)
  )
  const addableComponents = availableComponents.filter((name) => !attachedNames.includes(name))

  return (
    <div className="flex flex-1 flex-col overflow-hidden">
      <div className="border-border grid grid-cols-2 gap-2 border-b p-3">
        <TextField
          // Remounted per object so the inputs always start from the object
          // actually selected, rather than carrying over a previous edit.
          key={`${selected.id}:name:${selected.name}`}
          label="Name"
          value={selected.name}
          placeholder="GameObject"
          onCommit={(name) => updateMetadata(selected.id, { name })}
        />
        <TextField
          key={`${selected.id}:group:${selected.group}`}
          label="Group"
          value={selected.group}
          placeholder="None"
          onCommit={(group) => updateMetadata(selected.id, { group })}
        />
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
                isDragging={draggedComponent === name}
                onDragStart={setDraggedComponent}
                onDragEnd={() => setDraggedComponent(null)}
                onDrop={(sourceName, targetName) =>
                  reorderComponents(selected.id, sourceName, targetName)
                }
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

// One slot, since only one of Run/Stop is ever meaningful. The editor starts
// nothing here, so what this tracks is a request's progress, not a process.
function RunButton({
  canRunGames,
  activeProject,
  status,
  onRun,
  onStop
}: {
  canRunGames: boolean
  activeProject: ProjectSummary | null
  status: RunStatus
  onRun: (project: ProjectSummary) => void
  onStop: () => void
}): React.JSX.Element {
  const isBusy = status === 'starting' || status === 'stopping'

  const button =
    status === 'idle' ? (
      <Button
        size="sm"
        disabled={!activeProject || !canRunGames}
        onClick={() => activeProject && onRun(activeProject)}
      >
        <Play aria-hidden />
        Run
      </Button>
    ) : (
      <Button
        size="sm"
        variant={status === 'running' ? 'destructive' : 'secondary'}
        disabled={isBusy}
        onClick={onStop}
      >
        {isBusy ? <Loader2 className="animate-spin" aria-hidden /> : <Square aria-hidden />}
        {status === 'starting' ? 'Starting…' : status === 'stopping' ? 'Stopping…' : 'Stop'}
      </Button>
    )

  // Explain rather than leave a dead button: the reason is how the app was
  // launched, which isn't visible from in here.
  if (!canRunGames) {
    return (
      <Tooltip>
        <TooltipTrigger asChild>
          <span>{button}</span>
        </TooltipTrigger>
        <TooltipContent>{ORCHESTRATOR_UNAVAILABLE_MESSAGE}</TooltipContent>
      </Tooltip>
    )
  }

  return button
}

// The object actions sit directly above the canvas rather than in a panel of
// their own: they act on what the canvas shows, so they belong next to it.
function WorkspaceToolbar(): React.JSX.Element {
  const selectedObjectId = useEditorStore((state) => state.selectedObjectId)
  const addGameobject = useEditorStore((state) => state.addGameobject)
  const duplicateGameobject = useEditorStore((state) => state.duplicateGameobject)
  const deleteGameobject = useEditorStore((state) => state.deleteGameobject)
  const activeProject = useProjectStore((state) => state.activeProject)

  const runStatus = useRunStore((state) => state.status)
  const runGame = useRunStore((state) => state.run)
  const stopGame = useRunStore((state) => state.stop)

  const [saveState, setSaveState] = useState<'idle' | 'saving' | 'saved'>('idle')

  const canRunGames = isOrchestratorAvailable()

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

  return (
    <div className="border-border flex h-11 shrink-0 items-center gap-2 border-b px-3">
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

      <Button
        variant="outline"
        size="sm"
        disabled={!selectedObjectId}
        onClick={() => selectedObjectId && duplicateGameobject(selectedObjectId)}
      >
        <Copy aria-hidden />
        Duplicate
      </Button>

      {/* Save sits with Run rather than with the object actions: both act on
          the project as a whole, where the others act on one object. */}
      <div className="ml-auto flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={handleSave}
          disabled={!activeProject || saveState === 'saving'}
        >
          <Save aria-hidden />
          {saveState === 'saving' ? 'Saving…' : saveState === 'saved' ? 'Saved' : 'Save'}
        </Button>

        <RunButton
          canRunGames={canRunGames}
          activeProject={activeProject}
          status={runStatus}
          onRun={runGame}
          onStop={stopGame}
        />
      </div>
    </div>
  )
}

// Zoom bounds and step for the canvas. Wide enough to frame a whole scene or
// work on one object closely, without letting the screen rect vanish entirely.
const MIN_ZOOM = 0.1
const MAX_ZOOM = 4
const ZOOM_STEP = 1.1

interface CanvasView {
  x: number
  y: number
  scale: number
}

function clampZoom(scale: number): number {
  return Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, scale))
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
  // The canvas is infinite, so the view is panned and zoomed freely rather
  // than recomputed from the container: null until the first measure, which
  // is what frames the game screen on open.
  const [view, setView] = useState<CanvasView | null>(null)
  // Until the view has been moved by hand it keeps re-framing itself as the
  // container resizes. The first measure can arrive before the layout has
  // settled, and re-fitting is more useful than honouring that stale size.
  const hasAdjustedView = useRef(false)

  const fitToScreen = useCallback((width: number, height: number) => {
    const { scale, offsetX, offsetY } = computeFrameLayout(width, height)
    setView({ x: offsetX, y: offsetY, scale })
  }, [])

  useEffect(() => {
    const el = containerRef.current
    if (!el) return

    const observer = new ResizeObserver(([entry]) => {
      const { width, height } = entry.contentRect
      setSize({ width, height })
      if (!hasAdjustedView.current) fitToScreen(width, height)
    })
    observer.observe(el)
    return () => observer.disconnect()
  }, [fitToScreen])

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
    // Clicked the empty canvas or the game screen itself, not an object.
    if (e.target === e.target.getStage() || e.target.name() === SCREEN_NODE_NAME) {
      selectObject(null)
    }
  }

  // Zooms toward the pointer rather than the canvas origin, so the thing under
  // the cursor stays under the cursor.
  function handleWheel(e: Konva.KonvaEventObject<WheelEvent>): void {
    e.evt.preventDefault()
    const stage = stageRef.current
    const pointer = stage?.getPointerPosition()
    if (!stage || !pointer || !view) return

    const nextScale = clampZoom(view.scale * (e.evt.deltaY > 0 ? 1 / ZOOM_STEP : ZOOM_STEP))
    if (nextScale === view.scale) return

    hasAdjustedView.current = true

    const worldX = (pointer.x - view.x) / view.scale
    const worldY = (pointer.y - view.y) / view.scale
    setView({
      scale: nextScale,
      x: pointer.x - worldX * nextScale,
      y: pointer.y - worldY * nextScale
    })
  }

  const hasSize = size.width > 0 && size.height > 0
  const scale = view?.scale ?? 1

  return (
    <div ref={containerRef} className="bg-workspace relative min-h-0 flex-1 overflow-hidden">
      <h2 className="sr-only">Workspace</h2>
      {hasSize && view && (
        <>
          <span
            className="text-muted-foreground pointer-events-none absolute z-10 text-xs"
            style={{ left: view.x, top: view.y - 22 }}
          >
            Main Screen · {SCREEN_WIDTH}×{SCREEN_HEIGHT}
          </span>

          <div className="absolute right-3 bottom-3 z-10 flex items-center gap-1.5">
            <span className="text-muted-foreground text-xs tabular-nums">
              {Math.round(view.scale * 100)}%
            </span>
            <Button
              variant="secondary"
              size="sm"
              onClick={() => fitToScreen(size.width, size.height)}
            >
              <Maximize2 aria-hidden />
              Fit
            </Button>
          </div>

          <Stage
            ref={stageRef}
            width={size.width}
            height={size.height}
            x={view.x}
            y={view.y}
            scaleX={view.scale}
            scaleY={view.scale}
            draggable
            onDragEnd={(e) => {
              // Only the stage itself panning moves the view; a shape's own
              // drag bubbles up here too and must not be mistaken for a pan.
              if (e.target !== e.target.getStage()) return
              hasAdjustedView.current = true
              setView({ ...view, x: e.target.x(), y: e.target.y() })
            }}
            onWheel={handleWheel}
            onMouseDown={handleStageMouseDown}
          >
            <Layer>
              <Rect
                name={SCREEN_NODE_NAME}
                x={0}
                y={0}
                width={SCREEN_WIDTH}
                height={SCREEN_HEIGHT}
                fill={FRAME_FILL}
                stroke={SCREEN_STROKE}
                strokeWidth={1.5 / scale}
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
            <div className="flex h-full flex-col overflow-hidden">
              <WorkspaceToolbar />
              <Workspace />
            </div>
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

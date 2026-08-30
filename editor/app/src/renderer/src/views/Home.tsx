import { useEffect, useState } from 'react'
import { AlertCircle, FolderOpen, FolderPlus, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ScrollArea } from '@/components/ui/scroll-area'
import { cn } from '@/lib/utils'
import { useProjectStore } from '@/store/projectStore'
import { ProjectSummary } from '../../../shared/project'

const RELATIVE_TIME = new Intl.RelativeTimeFormat(undefined, { numeric: 'auto' })
const RELATIVE_UNITS: [Intl.RelativeTimeFormatUnit, number][] = [
  ['year', 60 * 60 * 24 * 365],
  ['month', 60 * 60 * 24 * 30],
  ['week', 60 * 60 * 24 * 7],
  ['day', 60 * 60 * 24],
  ['hour', 60 * 60],
  ['minute', 60]
]

function formatRelativeTime(iso: string): string {
  const seconds = Math.round((new Date(iso).getTime() - Date.now()) / 1000)
  for (const [unit, unitSeconds] of RELATIVE_UNITS) {
    if (Math.abs(seconds) >= unitSeconds) {
      return RELATIVE_TIME.format(Math.round(seconds / unitSeconds), unit)
    }
  }
  return RELATIVE_TIME.format(seconds, 'second')
}

function RowSkeleton(): React.JSX.Element {
  return (
    <div className="flex items-center gap-3 px-3 py-2.5">
      <div className="bg-muted size-4 shrink-0 animate-pulse rounded" />
      <div className="min-w-0 flex-1 space-y-1.5">
        <div className="bg-muted h-3.5 w-32 animate-pulse rounded" />
        <div className="bg-muted h-3 w-48 animate-pulse rounded" />
      </div>
    </div>
  )
}

function ProjectRow({ project }: { project: ProjectSummary }): React.JSX.Element {
  const openProject = useProjectStore((state) => state.openProject)
  const isLoading = useProjectStore((state) => state.isLoading)

  return (
    <button
      type="button"
      disabled={isLoading}
      onClick={() => openProject(project)}
      className={cn(
        'group flex w-full items-center gap-3 rounded-md px-3 py-2.5 text-left',
        'hover:bg-accent transition-colors duration-100',
        'focus-visible:ring-ring focus-visible:ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
        'disabled:pointer-events-none disabled:opacity-50'
      )}
    >
      <FolderOpen
        className="text-muted-foreground group-hover:text-foreground size-4 shrink-0 transition-colors duration-100"
        aria-hidden
      />
      <div className="min-w-0 flex-1">
        <p className="truncate text-sm font-medium">{project.name}</p>
        <p className="text-muted-foreground truncate text-xs">{project.path}</p>
      </div>
      <span className="text-muted-foreground shrink-0 text-xs">
        {formatRelativeTime(project.lastOpenedAt)}
      </span>
    </button>
  )
}

function NewProjectDialog(): React.JSX.Element {
  const [name, setName] = useState('')
  const open = useProjectStore((state) => state.isNewProjectDialogOpen)
  const setOpen = useProjectStore((state) => state.setNewProjectDialogOpen)
  const isLoading = useProjectStore((state) => state.isLoading)
  const error = useProjectStore((state) => state.error)
  const createProject = useProjectStore((state) => state.createProject)

  async function handleCreate(): Promise<void> {
    if (!name.trim()) return
    const created = await createProject(name)
    if (created) {
      setName('')
      setOpen(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus aria-hidden />
          New Project
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New project</DialogTitle>
          <DialogDescription>
            Creates a new folder under your LevelCraft Projects directory.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-2">
          <Label htmlFor="new-project-name">Project name</Label>
          <Input
            id="new-project-name"
            autoFocus
            placeholder="My Game"
            value={name}
            onChange={(e) => setName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') handleCreate()
            }}
            disabled={isLoading}
          />
          {error && (
            <p className="text-destructive flex items-center gap-1.5 text-xs">
              <AlertCircle className="size-3.5 shrink-0" aria-hidden />
              {error}
            </p>
          )}
        </div>

        <DialogFooter>
          <DialogClose asChild>
            <Button variant="ghost" disabled={isLoading}>
              Cancel
            </Button>
          </DialogClose>
          <Button onClick={handleCreate} disabled={isLoading || !name.trim()}>
            Create
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function EmptyState(): React.JSX.Element {
  return (
    <div className="flex flex-col items-center justify-center gap-3 py-16 text-center">
      <FolderPlus className="text-muted-foreground size-8" aria-hidden />
      <div className="space-y-1">
        <p className="text-sm font-medium">No projects yet</p>
        <p className="text-muted-foreground text-xs">Create a project to start building a scene.</p>
      </div>
    </div>
  )
}

function ErrorState({
  message,
  onRetry
}: {
  message: string
  onRetry: () => void
}): React.JSX.Element {
  return (
    <div className="border-destructive/20 bg-destructive/5 flex flex-col items-start gap-3 rounded-lg border p-4">
      <div className="text-destructive flex items-center gap-2">
        <AlertCircle className="size-4" aria-hidden />
        <p className="text-sm font-medium">Couldn&apos;t load your projects</p>
      </div>
      <p className="text-muted-foreground text-xs">{message}</p>
      <Button variant="secondary" size="sm" onClick={onRetry}>
        Try again
      </Button>
    </div>
  )
}

function Home(): React.JSX.Element {
  const projects = useProjectStore((state) => state.projects)
  const isLoading = useProjectStore((state) => state.isLoading)
  const listError = useProjectStore((state) => state.listError)
  const error = useProjectStore((state) => state.error)
  const backendReachable = useProjectStore((state) => state.backendReachable)
  const checkBackend = useProjectStore((state) => state.checkBackend)
  const refreshProjects = useProjectStore((state) => state.refreshProjects)
  const openProjectFromPath = useProjectStore((state) => state.openProjectFromPath)

  useEffect(() => {
    refreshProjects()
    checkBackend()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  async function handleBrowse(): Promise<void> {
    const path = await window.api.project.browseForFolder()
    if (path) await openProjectFromPath(path)
  }

  // Nothing else to show when the list itself never loaded — the big
  // centered error state applies. If we already have a list, an error
  // (e.g. a project failed to open) shows as a compact inline banner
  // instead, so it doesn't blow away projects the user can still see.
  const listFailedToLoad = !isLoading && projects.length === 0 && listError

  return (
    <div className="bg-background text-foreground flex h-screen flex-col">
      <header
        className={cn(
          'flex h-11 shrink-0 items-center justify-between border-b pr-3',
          window.api.platform === 'darwin' ? 'pl-20' : 'pl-3'
        )}
        style={{ WebkitAppRegion: 'drag' } as React.CSSProperties}
      >
        <span className="text-sm font-semibold tracking-tight">LevelCraft</span>
        <div className="flex gap-2" style={{ WebkitAppRegion: 'no-drag' } as React.CSSProperties}>
          <Button variant="outline" size="sm" onClick={handleBrowse} disabled={isLoading}>
            Open Existing
          </Button>
          <NewProjectDialog />
        </div>
      </header>

      <div className="mx-auto flex min-h-0 w-full max-w-2xl flex-1 flex-col px-6 py-8">
        <h1 className="mb-1 shrink-0 text-sm font-medium">Projects</h1>
        <p className="text-muted-foreground mb-4 shrink-0 text-xs">
          Open a recent project, or create a new one to get started.
        </p>

        {backendReachable === false && (
          <p className="text-destructive mb-4 flex shrink-0 items-center gap-1.5 text-xs">
            <AlertCircle className="size-3.5 shrink-0" aria-hidden />
            Can&apos;t reach the LevelCraft backend — make sure editor/backend is running.
          </p>
        )}

        {listFailedToLoad ? (
          <ErrorState message={listError} onRetry={refreshProjects} />
        ) : (
          <div className="flex min-h-0 flex-1 flex-col">
            {error && projects.length > 0 && (
              <p className="text-destructive mb-3 flex shrink-0 items-center gap-1.5 text-xs">
                <AlertCircle className="size-3.5 shrink-0" aria-hidden />
                {error}
              </p>
            )}
            {isLoading && projects.length === 0 ? (
              <div className="space-y-1">
                <RowSkeleton />
                <RowSkeleton />
                <RowSkeleton />
              </div>
            ) : projects.length === 0 ? (
              <div className="flex flex-1 items-center justify-center">
                <EmptyState />
              </div>
            ) : (
              <ScrollArea className="min-h-0 flex-1">
                <div className="space-y-0.5">
                  {projects.map((project) => (
                    <ProjectRow key={project.id} project={project} />
                  ))}
                </div>
              </ScrollArea>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

export default Home

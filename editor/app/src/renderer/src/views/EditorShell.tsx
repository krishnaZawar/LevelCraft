import { useEffect } from 'react'
import { LayoutGrid, ListTree, SlidersHorizontal, X } from 'lucide-react'
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from '@/components/ui/resizable'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { useEditorStore } from '@/store/editorStore'
import { useProjectStore } from '@/store/projectStore'

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
      <span className="text-sm font-medium">{activeProject?.name}</span>
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

function Workspace(): React.JSX.Element {
  return (
    <div className="flex h-full flex-col overflow-hidden">
      <h2 className="sr-only">Workspace</h2>
      <div
        className="flex-1"
        style={{
          backgroundImage:
            'radial-gradient(color-mix(in oklch, var(--muted-foreground) 25%, transparent) 1px, transparent 1px)',
          backgroundSize: '16px 16px'
        }}
      />
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
              <div className="flex-1 overflow-y-auto p-3" />
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
                  <div className="flex-1 overflow-y-auto p-3" />
                </div>
              </ResizablePanel>
            </ResizablePanelGroup>
          </ResizablePanel>

          <ResizableHandle withHandle />

          <ResizablePanel defaultSize="20%" minSize="15%" maxSize="35%">
            <div className="flex h-full flex-col overflow-hidden">
              <PanelHeading icon={SlidersHorizontal}>Attributes</PanelHeading>
              <div className="flex-1 overflow-y-auto p-3" />
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  )
}

export default EditorShell

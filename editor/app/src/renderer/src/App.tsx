import { TooltipProvider } from '@/components/ui/tooltip'

function App(): React.JSX.Element {
  return (
    <TooltipProvider delayDuration={200}>
      <div className="bg-background text-foreground grid h-screen grid-cols-[240px_1fr_280px] grid-rows-[1fr_140px]">
        <aside className="border-border row-span-2 overflow-y-auto border-r p-3">
          <h2 className="text-muted-foreground mb-2 text-xs font-semibold tracking-wide uppercase">
            Hierarchy
          </h2>
        </aside>

        <main className="border-border overflow-hidden border-b">
          <h2 className="sr-only">Workspace</h2>
        </main>

        <aside className="border-border row-span-2 overflow-y-auto border-l p-3">
          <h2 className="text-muted-foreground mb-2 text-xs font-semibold tracking-wide uppercase">
            Attributes
          </h2>
        </aside>

        <section className="overflow-y-auto p-3">
          <h2 className="text-muted-foreground mb-2 text-xs font-semibold tracking-wide uppercase">
            Utility
          </h2>
        </section>
      </div>
    </TooltipProvider>
  )
}

export default App

function App(): React.JSX.Element {
  return (
    <div className="grid h-screen grid-cols-[240px_1fr_280px] grid-rows-[1fr_140px] bg-neutral-950 text-neutral-200">
      <aside className="row-span-2 overflow-y-auto border-r border-neutral-800 p-3">
        <h2 className="mb-2 text-xs font-semibold tracking-wide text-neutral-400 uppercase">
          Hierarchy
        </h2>
      </aside>

      <main className="overflow-hidden border-b border-neutral-800">
        <h2 className="sr-only">Workspace</h2>
      </main>

      <aside className="row-span-2 overflow-y-auto border-l border-neutral-800 p-3">
        <h2 className="mb-2 text-xs font-semibold tracking-wide text-neutral-400 uppercase">
          Attributes
        </h2>
      </aside>

      <section className="overflow-y-auto p-3">
        <h2 className="mb-2 text-xs font-semibold tracking-wide text-neutral-400 uppercase">
          Utility
        </h2>
      </section>
    </div>
  )
}

export default App

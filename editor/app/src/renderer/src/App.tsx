import { useEffect } from 'react'
import { TooltipProvider } from '@/components/ui/tooltip'
import { useProjectStore } from '@/store/projectStore'
import EditorShell from '@/views/EditorShell'
import Home from '@/views/Home'

function App(): React.JSX.Element {
  const activeProject = useProjectStore((state) => state.activeProject)
  const isProjectOpen = Boolean(activeProject)

  useEffect(() => {
    if (isProjectOpen) {
      window.api.window.maximize()
    } else {
      window.api.window.unmaximize()
    }
    window.api.menu.notifyProjectOpen(isProjectOpen)
  }, [isProjectOpen])

  return (
    <TooltipProvider delayDuration={200}>
      {isProjectOpen ? <EditorShell /> : <Home />}
    </TooltipProvider>
  )
}

export default App

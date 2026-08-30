import { useEffect } from 'react'
import { TooltipProvider } from '@/components/ui/tooltip'
import { useEditorStore } from '@/store/editorStore'
import { useProjectStore } from '@/store/projectStore'
import EditorShell from '@/views/EditorShell'
import Home from '@/views/Home'

function App(): React.JSX.Element {
  const activeProject = useProjectStore((state) => state.activeProject)
  const isProjectOpen = Boolean(activeProject)
  const resetEditorState = useEditorStore((state) => state.reset)

  useEffect(() => {
    if (isProjectOpen) {
      window.api.window.maximize()
    } else {
      window.api.window.unmaximize()
      // Clear any previous project's scene state so it can't linger and
      // briefly flash once a different project is opened next.
      resetEditorState()
    }
    window.api.menu.notifyProjectOpen(isProjectOpen)
  }, [isProjectOpen, resetEditorState])

  return (
    <TooltipProvider delayDuration={200}>
      {isProjectOpen ? <EditorShell /> : <Home />}
    </TooltipProvider>
  )
}

export default App

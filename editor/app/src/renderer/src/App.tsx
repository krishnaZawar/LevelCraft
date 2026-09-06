import { useEffect } from 'react'
import { Toaster } from '@/components/ui/sonner'
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
      // Closing the project while the Play window is open would otherwise
      // leave it (and builder/backend) running with no way back to it.
      window.api.builder.stop()
    }
    window.api.menu.notifyProjectOpen(isProjectOpen)
  }, [isProjectOpen, resetEditorState])

  return (
    <TooltipProvider delayDuration={200}>
      {isProjectOpen ? <EditorShell /> : <Home />}
      <Toaster position="bottom-right" />
    </TooltipProvider>
  )
}

export default App

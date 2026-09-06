import PlayView from './views/PlayView'

// Root component for the separate Play window (see main/index.ts's
// createPlayWindow) — loaded from the same bundle as the main editor via
// the '#play' URL hash (see main.tsx), not a second electron-vite entry.
function PlayWindowApp(): React.JSX.Element {
  return <PlayView onStop={() => window.api.builder.stop()} />
}

export default PlayWindowApp

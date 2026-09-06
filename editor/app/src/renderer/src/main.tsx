import './assets/main.css'

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import PlayWindowApp from './PlayWindowApp'
import { ErrorBoundary } from './ErrorBoundary'
import { initMenuBridge } from './menuBridge'

// The Play window (see main/index.ts's createPlayWindow) loads this same
// bundle with a '#play' hash instead of getting its own electron-vite
// entry point — cheaper than a second renderer target for one small window.
const isPlayWindow = window.location.hash === '#play'

if (!isPlayWindow) {
  initMenuBridge()
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ErrorBoundary>{isPlayWindow ? <PlayWindowApp /> : <App />}</ErrorBoundary>
  </StrictMode>
)

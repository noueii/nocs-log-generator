import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import './index.css'
import { router } from './router'
import { ErrorBoundary } from './components/ErrorBoundary'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ErrorBoundary level="critical" showDetails={import.meta.env.DEV}>
      <RouterProvider router={router} />
    </ErrorBoundary>
  </StrictMode>,
)

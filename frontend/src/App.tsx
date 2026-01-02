import React from 'react'
import { Router } from './components/Router/Router'
import { AppRoutes } from './components/Router/AppRoutes'
import { StoreProvider } from './store/StoreContext'

export const App = () => {
  return (
    <StoreProvider>
      <Router>
        <AppRoutes />
      </Router>
    </StoreProvider>
  )
}

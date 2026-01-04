import React, { useEffect } from 'react'
import { BrowserRouter, useLocation, useNavigate, useNavigationType } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useStore } from '../../store'

interface RouterProps {
  children: React.ReactNode
}

const RouterObserver: React.FC = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const navigationType = useNavigationType()
  const { route } = useStore()

  useEffect(() => {
    route.setNavigateFunction(navigate)
  }, [route, navigate])

  useEffect(() => {
    if (navigationType !== 'POP') {
      const previousPath = route.currentPath

      route.routeState.currentPath = location.pathname

      if (previousPath !== location.pathname) {
        route.routeState.previousPath = previousPath

        const lastHistoryPath =
          route.routeState.navigationHistory[route.routeState.navigationHistory.length - 1]
        if (lastHistoryPath !== location.pathname) {
          route.routeState.navigationHistory.push(location.pathname)
        }
      }

      route.updateTelegramBackButton()
    } else {
      route.routeState.currentPath = location.pathname

      if (route.routeState.navigationHistory.length > 1) {
        route.routeState.navigationHistory.pop()

        const newHistoryLength = route.routeState.navigationHistory.length
        route.routeState.previousPath =
          newHistoryLength > 1 ? route.routeState.navigationHistory[newHistoryLength - 2] : null
      }

      route.updateTelegramBackButton()
    }
  }, [location, navigationType, route])

  return null
}

export const Router: React.FC<RouterProps> = observer(({ children }) => {
  return (
    <BrowserRouter>
      <RouterObserver />
      {children}
    </BrowserRouter>
  )
})

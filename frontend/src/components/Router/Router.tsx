import React, { useEffect } from 'react'
import { BrowserRouter, useLocation, useNavigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useStore } from '../../store'

interface RouterProps {
  children: React.ReactNode
}

const RouterObserver: React.FC = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const { route } = useStore()

  useEffect(() => {
    route.setNavigateFunction(navigate)
  }, [route, navigate])

  useEffect(() => {
    if (route.currentPath !== location.pathname) {
      route.navigate(location.pathname)
    }
  }, [location.pathname, route])

  useEffect(() => {
    const handlePopState = () => {
      route.goBack()
    }

    window.addEventListener('popstate', handlePopState)
    return () => {
      window.removeEventListener('popstate', handlePopState)
    }
  }, [route])

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

import { Loader } from '@/components/Loader/Loader'
import { observer } from 'mobx-react-lite'
import React, { useEffect } from 'react'
import { Route, Routes } from 'react-router-dom'
import { TabBar } from '../TabBar/TabBar'
import { ErrorPage } from '@/pages/Error/ErrorPage'
import { postEvent } from '@telegram-apps/sdk'

const HomePage = React.lazy(() =>
  import('@/pages/Home/HomePage').then(module => ({ default: module.HomePage }))
)
const BrowsePage = React.lazy(() =>
  import('@/pages/Browse/BrowsePage').then(module => ({ default: module.BrowsePage }))
)
const SettingsPage = React.lazy(() =>
  import('@/pages/Settings/SettingsPage').then(module => ({ default: module.SettingsPage }))
)

export const AppRoutes = observer(() => {
  useEffect(() => {
    if (window.Telegram?.WebApp.initDataUnsafe?.user) {
      postEvent('web_app_request_fullscreen')
    }
  }, [])

  return (
    <React.Suspense fallback={<Loader />}>
      <Routes>
        <Route path='/' element={<HomePage />} />
        <Route path='/browse' element={<BrowsePage />} />
        <Route path='/settings' element={<SettingsPage />} />
        <Route path='*' element={<ErrorPage />} />
      </Routes>
      <TabBar />
    </React.Suspense>
  )
})

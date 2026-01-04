import { Loader } from '@/components/Loader/Loader'
import { ErrorPage } from '@/pages/Error/ErrorPage'
import { postEvent } from '@telegram-apps/sdk'
import { observer } from 'mobx-react-lite'
import React, { useEffect } from 'react'
import { Route, Routes } from 'react-router-dom'
import { NavBar } from '../NavBar/NavBar'
import { TabBar } from '../TabBar/TabBar'

const HomePage = React.lazy(() =>
  import('@/pages/Home/HomePage').then(module => ({ default: module.HomePage }))
)
const BrowsePage = React.lazy(() =>
  import('@/pages/Browse/BrowsePage').then(module => ({ default: module.BrowsePage }))
)
const CreatePage = React.lazy(() =>
  import('@/pages/Create/CreatePage').then(module => ({ default: module.CreatePage }))
)
const ProfilePage = React.lazy(() =>
  import('@/pages/Profile/ProfilePage').then(module => ({ default: module.ProfilePage }))
)

export const AppRoutes = observer(() => {
  useEffect(() => {
    if (window.Telegram?.WebApp.initDataUnsafe?.user) {
      postEvent('web_app_request_fullscreen')
    }
  }, [])

  return (
    <React.Suspense fallback={<Loader />}>
      <NavBar />
      <Routes>
        <Route path='/' element={<HomePage />} />
        <Route path='/browse' element={<BrowsePage />} />
        <Route path='/create' element={<CreatePage />} />
        <Route path='/profile' element={<ProfilePage />} />
        <Route path='*' element={<ErrorPage />} />
      </Routes>
      <TabBar />
    </React.Suspense>
  )
})

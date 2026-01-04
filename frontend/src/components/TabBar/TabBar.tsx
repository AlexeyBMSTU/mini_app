import { SyntheticEvent, useMemo } from 'react'
import { observer } from 'mobx-react-lite'
import { useStore } from '@/store'
import AccountCircleRoundedIcon from '@mui/icons-material/AccountCircleRounded'
import Tab from '@mui/material/Tab'
import Tabs from '@mui/material/Tabs'
import AddHomeWorkOutlinedIcon from '@mui/icons-material/AddHomeWorkOutlined'
import MapsHomeWorkOutlinedIcon from '@mui/icons-material/MapsHomeWorkOutlined'
import MapsHomeWorkRoundedIcon from '@mui/icons-material/MapsHomeWorkRounded'
import AccountCircleOutlinedIcon from '@mui/icons-material/AccountCircleOutlined'
import AddHomeWorkRoundedIcon from '@mui/icons-material/AddHomeWorkRounded'

import styles from './TabBar.module.css'

const TAB_ROUTES = ['/browse', '/create', '/profile'] as const

export const TabBar = observer(() => {
  const { route } = useStore()

  const activeTab = useMemo(() => {
    const currentPath = route.currentPath
    const index = TAB_ROUTES.indexOf(currentPath as (typeof TAB_ROUTES)[number])
    return index !== -1 ? index : 0
  }, [route.currentPath])

  if (route.currentPath === '/') {
    return null
  }

  const handleRoute = (_event: SyntheticEvent, newValue: number) => {
    const routePath = TAB_ROUTES[newValue]
    if (routePath) {
      route.navigate(routePath)
    }
  }

  return (
    <div className={styles.tabBar} style={{ backgroundColor: 'white' }}>
      <Tabs
        value={activeTab}
        onChange={handleRoute}
        aria-label='Navigation tabs'
        slotProps={{
          indicator: {
            style: { top: 0 },
          },
          list: {
            style: { gap: '15px' },
          },
        }}
      >
        <Tab
          icon={
            activeTab === 0 ? (
              <MapsHomeWorkRoundedIcon fontSize='large' color='primary' />
            ) : (
              <MapsHomeWorkOutlinedIcon fontSize='large' color='primary' />
            )
          }
          aria-label='Houmes'
        />
        <Tab
          icon={
            activeTab === 1 ? (
              <AddHomeWorkRoundedIcon fontSize='large' color='primary' />
            ) : (
              <AddHomeWorkOutlinedIcon fontSize='large' color='primary' />
            )
          }
          aria-label='Create'
        />
        <Tab
          icon={
            activeTab === 2 ? (
              <AccountCircleRoundedIcon fontSize='large' color='primary' />
            ) : (
              <AccountCircleOutlinedIcon fontSize='large' color='primary' />
            )
          }
          aria-label='Profile'
        />
      </Tabs>
    </div>
  )
})

import { SyntheticEvent, useMemo } from 'react'
import { observer } from 'mobx-react-lite'
import { useStore } from '@/store'
import AccountCircleRoundedIcon from '@mui/icons-material/AccountCircleRounded'
import FormatListBulletedRoundedIcon from '@mui/icons-material/FormatListBulletedRounded'
import Tab from '@mui/material/Tab'
import Tabs from '@mui/material/Tabs'
import styles from './TabBar.module.css'

const TAB_ROUTES = ['/chats', '/settings'] as const

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
    <Tabs
      value={activeTab}
      onChange={handleRoute}
      className={styles.tabBar}
      aria-label='Navigation tabs'
    >
      <Tab icon={<FormatListBulletedRoundedIcon color='primary' />} aria-label='Chats' />
      <Tab icon={<AccountCircleRoundedIcon color='primary' />} aria-label='Settings' />
    </Tabs>
  )
})

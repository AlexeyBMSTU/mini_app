import { SyntheticEvent, useMemo } from 'react'
import { observer } from 'mobx-react-lite'
import { useStore } from '@/store'
import AccountCircleRoundedIcon from '@mui/icons-material/AccountCircleRounded'
import Tab from '@mui/material/Tab'
import Tabs from '@mui/material/Tabs'
import HomeOutlinedIcon from '@mui/icons-material/HomeOutlined';
import HomeRoundedIcon from '@mui/icons-material/HomeRounded';
import AccountCircleOutlinedIcon from '@mui/icons-material/AccountCircleOutlined';
import styles from './TabBar.module.css'

const TAB_ROUTES = ['/browse', '/settings'] as const

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
    <div
      className={styles.tabBar}
      style={{backgroundColor: 'white'}}
    >
      <Tabs
        value={activeTab}
        onChange={handleRoute}
        aria-label='Navigation tabs'
        sx={{paddingBottom: '3px'}}
        slotProps={{
          indicator: {
            style: { display: 'none' }
          }
        }}
      >
        <Tab icon={activeTab === 0 ? <HomeRoundedIcon  fontSize='large'  color='primary' /> : <HomeOutlinedIcon fontSize='large' 
        color='primary' />} aria-label='Houmes' />
        <Tab icon={activeTab === 1 ? <AccountCircleRoundedIcon fontSize='large' color='primary' /> : <AccountCircleOutlinedIcon fontSize='large' color='primary' />} aria-label='Profile' />
      </Tabs>
    </div>
  )
})

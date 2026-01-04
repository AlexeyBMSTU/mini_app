import { ReactNode } from 'react'
import styles from './PurePage.module.css'
import { observer } from 'mobx-react-lite'
import { Box } from '@mui/material'

export const PurePage = observer(({ children }: { children: ReactNode }) => {
  return (
    <div className={styles.container}>
      <Box
        sx={{
          width: '100%',
          maxWidth: '768px',
          display: 'flex',
          justifyContent: 'center',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        {children}
      </Box>
    </div>
  )
})

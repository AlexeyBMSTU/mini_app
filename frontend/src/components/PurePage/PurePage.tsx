import { ReactNode } from 'react'
import styles from './PurePage.module.css'
import { observer } from 'mobx-react-lite'

export const PurePage = observer(({ children }: { children: ReactNode }) => {
  return <div className={styles.container}>{children}</div>
})

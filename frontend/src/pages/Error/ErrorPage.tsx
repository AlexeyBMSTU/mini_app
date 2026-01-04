import { observer } from 'mobx-react-lite'
import ErrorIcon from '@mui/icons-material/Error'
import styles from './ErrorPage.module.css'
import { PurePage } from '@/components/PurePage/PurePage'

export const ErrorPage = observer(() => {
  return (
    <PurePage>
      <section className={styles.root}>
        <ErrorIcon sx={{ fontSize: 60 }} color='error' />
        <p style={{ fontSize: 24 }}>Упс... Произошла ошибка</p>
      </section>
    </PurePage>
  )
})

import { PurePage } from '@/components/PurePage/PurePage'
import { useTelegram } from '@/hooks/useTelegram'
import styles from '@/pages/Home/HomePage.module.css'
import { useStore } from '@/store'
import { TypingEffect } from '@/utils/TypingEffect'
import { ArrowForwardRounded } from '@mui/icons-material'
import Button from '@mui/material/Button'
import Checkbox from '@mui/material/Checkbox'
import FormControlLabel from '@mui/material/FormControlLabel'
import FormGroup from '@mui/material/FormGroup'
import { observer } from 'mobx-react-lite'
import { useEffect, useState } from 'react'

export const HomePage = observer(() => {
  const [disabled, setDisabled] = useState(true)
  const [navigating, setNavigating] = useState(false)
  const { route } = useStore()
  const { user } = useTelegram()
  
  const handleContinue = () => {
    setNavigating(true);
    route.navigate('/chats');
  }
  useEffect(() => {
    if (user) {
      console.log('User ID из хука:', user.id)
    }
  }, [user])
  
  return (
    <PurePage>
      <section className={styles.welcomeSection}>
        <div style={{ color: '#1976d2' }} className={styles.welcomeText}>
          <TypingEffect isBold text='SERVATORY' />
        </div>
        <div style={{ color: '#1976d2', fontSize: 24 }} className={styles.welcomeText}>
          <TypingEffect text='Будь ближе к&nbsp;недвижимости' />
        </div>
        <FormGroup className={styles.welcomeButton}>
          <Button
            size='large'
            disabled={disabled || navigating}
            variant='outlined'
            onClick={handleContinue}
          >
            Продолжить
            <ArrowForwardRounded
              color={disabled || navigating ? 'disabled' : 'primary'}
              fontSize='large'
            />
          </Button>
          <FormControlLabel
            control={
              <Checkbox
                style={{ color: '#1976d2' }}
                onClick={() => setDisabled(disabled => !disabled)}
                disabled={navigating}
              />
            }
            label={<span style={{ color: '#1976d2' }}>Согласие на обработку перс. данных</span>}
          />
        </FormGroup>
      </section>
    </PurePage>
  )
})

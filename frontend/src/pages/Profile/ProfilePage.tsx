import { observer } from 'mobx-react-lite'
import { useTelegram } from '@/hooks/useTelegram'
import { PurePage } from '@/components/PurePage/PurePage'
import { Avatar, Card, CardContent, Container, Typography } from '@mui/material'
import { motion } from 'motion/react'
import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'
import ListItemText from '@mui/material/ListItemText'
import Divider from '@mui/material/Divider'
import styles from '@/pages/Profile/ProfilePage.module.css'
import { ArrowForwardRounded } from '@mui/icons-material'

export const ProfilePage = observer(() => {
  const { user } = useTelegram()

  if (!user) {
    return (
      <PurePage>
        <div className={styles.errorContainer}>
          <Typography variant='h6'>Не удалось загрузить данные пользователя</Typography>
        </div>
      </PurePage>
    )
  }

  return (
    <PurePage>
      <motion.div
        style={{ width: '100%' }}
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5 }}
      >
        <Card className={styles.profileCard} elevation={0}>
          <CardContent className={styles.profileContent}>
            <motion.div
              className={styles.avatarContainer}
              initial={{ scale: 0.8, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              transition={{ duration: 0.6, ease: 'easeOut' }}
            >
              <Avatar
                alt={`${user.first_name} ${user.last_name || ''}`}
                src={`https://t.me/i/userpic/320/${user.username}.jpg`}
                className={styles.avatar}
                sx={{
                  width: '120px',
                  height: '120px',
                  border: '4px solid #f0f0f0',
                  boxShadow: '0 4px 10px rgba(0, 0, 0, 0.1)',
                  transition: 'transform 0.3s ease',
                  '&:hover': {
                    transform: 'scale(1.05)',
                  },
                }}
              />
            </motion.div>

            <div className={styles.userInfo}>
              <Typography variant='h5' component='h1'>
                {user.first_name} {user.last_name || ''}
              </Typography>

              {user.username && (
                <Typography variant='subtitle1' color='textSecondary'>
                  @{user.username}
                </Typography>
              )}
            </div>
          </CardContent>
        </Card>
        <Container>
          <List>
            <ListItem>
              <ListItemText primary='Объявления' />
              <ListItemText
                sx={{ justifyContent: 'flex-end', display: 'flex', color: 'gray' }}
                primary='0'
              />
              <ArrowForwardRounded color='primary' fontSize='medium' sx={{ marginLeft: '5px' }} />
            </ListItem>
            <Divider component='li' />
            <ListItem>
              <ListItemText primary='Избранное' />
              <ArrowForwardRounded color='primary' fontSize='medium' />
            </ListItem>
            <Divider component='li' />
          </List>
        </Container>
      </motion.div>
    </PurePage>
  )
})

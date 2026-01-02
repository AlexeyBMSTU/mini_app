import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'
import Divider from '@mui/material/Divider'
import ListItemText from '@mui/material/ListItemText'
import ListItemAvatar from '@mui/material/ListItemAvatar'
import Avatar from '@mui/material/Avatar'
import { observer } from 'mobx-react-lite'
import { PurePage } from '@/components/PurePage/PurePage'
import styles from './ChatsPage.module.css'

export const ChatsPage = observer(() => {
  const map = {
    name: 'Алкаш с авито',
    message: 'Интересует предложение?',
  }
  return (
    <PurePage>
      <List sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}>
        <ListItem className={styles.item} alignItems='flex-start'>
          <ListItemAvatar>
            <Avatar sx={{ bgcolor: '#1976d2' }} alt={map.name} src='/static/images/avatar/1.jpg' />
          </ListItemAvatar>
          <ListItemText primary={map.name} secondary={map.message} />
        </ListItem>
        <Divider variant='inset' component='li' />
      </List>
    </PurePage>
  )
})

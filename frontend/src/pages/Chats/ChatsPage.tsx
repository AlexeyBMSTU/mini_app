import { Loader } from '@/components/Loader/Loader'
import { PurePage } from '@/components/PurePage/PurePage'
import { useStore } from '@/store/StoreContext'
import Avatar from '@mui/material/Avatar'
import Divider from '@mui/material/Divider'
import List from '@mui/material/List'
import ListItem from '@mui/material/ListItem'
import ListItemAvatar from '@mui/material/ListItemAvatar'
import ListItemText from '@mui/material/ListItemText'
import { observer } from 'mobx-react-lite'
import { useEffect } from 'react'
import { ErrorPage } from '../Error/ErrorPage'
import styles from './ChatsPage.module.css'

export const ChatsPage = observer(() => {
  const { chat } = useStore()

  useEffect(() => {
    if (chat.chatState.chats.length === 0) {
      chat.fetchChats()
    }
  }, [chat])

  if (chat.loading) {
    return (
      <PurePage>
        <Loader />
      </PurePage>
    )
  }

  if (chat.error) {
    return (
      <PurePage>
        <ErrorPage />
      </PurePage>
    )
  }
  return (
    <PurePage>
      <List sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}>
        {chat.chats.map(chatItem => (
          <div key={chatItem.id}>
            <ListItem className={styles.item} alignItems='flex-start'>
              <ListItemAvatar>
                <Avatar
                  sx={{ bgcolor: '#1976d2' }}
                  alt={chatItem.users[0].public_user_profile.avatar.default}
                  src={chatItem.users[0].public_user_profile.avatar.default}
                />
              </ListItemAvatar>
              <ListItemText
                primary={chatItem.users[0].name}
                secondary={
                  <>
                    <span>{chatItem.last_message.content.text}</span>
                    {chatItem.updated && (
                      <span style={{ fontSize: '0.75rem', color: '#999', marginLeft: '8px' }}>
                        {chatItem.updated}
                      </span>
                    )}
                  </>
                }
              />
            </ListItem>
            <Divider variant='inset' component='li' />
          </div>
        ))}
      </List>
    </PurePage>
  )
})

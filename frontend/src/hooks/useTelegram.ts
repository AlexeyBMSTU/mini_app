import { useState, useEffect } from 'react'
import { TelegramUser } from '../types/telegram'
import telegramService from '../services/telegramService'

export const useTelegram = () => {
  const [user, setUser] = useState<TelegramUser | null>(null)
  const [loading, setLoading] = useState(true)
  const [isTelegram, setIsTelegram] = useState(false)

  useEffect(() => {
    const initTelegram = () => {
      try {
        const telegramInitialized = telegramService.init()

        if (telegramInitialized) {
          setIsTelegram(true)
          const userData = telegramService.getUser()
          if (userData) {
            setUser(userData)
          }
        } else {
          setIsTelegram(false)
          const mockUser = telegramService.getMockUser()
          setUser(mockUser)
        }
      } catch (error) {
        console.error('Ошибка инициализации:', error)
        setIsTelegram(false)
      } finally {
        setLoading(false)
      }
    }

    initTelegram()
  }, [])

  const showAlert = (message: string) => {
    telegramService.showAlert(message)
  }

  const closeApp = () => {
    telegramService.close()
  }

  return {
    user,
    loading,
    isTelegram,
    showAlert,
    closeApp,
  }
}

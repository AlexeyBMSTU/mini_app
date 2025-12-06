import { useState, useEffect } from 'react';
import { TelegramUser } from '../types/telegram';
import telegramService from '../services/telegramService';

export const useTelegram = () => {
  const [user, setUser] = useState<TelegramUser | null>(null);
  const [loading, setLoading] = useState(true);
  const [isTelegram, setIsTelegram] = useState(false);

  useEffect(() => {
    const initTelegram = () => {
      try {
        const isLocalHost = window.location.hostname === 'localhost' || 
                           window.location.hostname === '127.0.0.1' || 
                           window.location.hostname.startsWith('192.168.');

        const telegramInitialized = telegramService.init();
        
        if (telegramInitialized && !isLocalHost) {
          setIsTelegram(true);
          const userData = telegramService.getUser();
          if (userData) {
            setUser(userData);
          }
          console.log('Запущено в Telegram');
        } else {
          setIsTelegram(false);
          console.log('Запущено в браузере. Используется режим разработки.');
          setUser(telegramService.getMockUser());
        }
      } catch (error) {
        console.error('Ошибка инициализации:', error);
        setIsTelegram(false);
      } finally {
        setLoading(false);
      }
    };

    initTelegram();
  }, []);

  const showAlert = (message: string) => {
    telegramService.showAlert(message);
  };

  const closeApp = () => {
    telegramService.close();
  };

  return {
    user,
    loading,
    isTelegram,
    showAlert,
    closeApp
  };
};
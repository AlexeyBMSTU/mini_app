import { useState, useEffect } from 'react';
import { TelegramUser } from '../types/telegram';
import telegramService from '../services/telegramService';
import apiService from '../services/apiService';

interface UseApiReturn {
  user: TelegramUser | null;
  loading: boolean;
  error: string | null;
  authenticate: () => Promise<void>;
  getUserData: () => Promise<any>;
  saveUserData: (data: any) => Promise<void>;
}

export const useApi = (): UseApiReturn => {
  const [user, setUser] = useState<TelegramUser | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Инициализация пользователя из Telegram
  useEffect(() => {
    const initUser = async () => {
      try {
        const telegramUser = telegramService.getUser();
        if (telegramUser) {
          setUser(telegramUser);
          // Автоматическая аутентификация при загрузке
          await authenticateUser(telegramUser);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to initialize user');
      } finally {
        setLoading(false);
      }
    };

    initUser();
  }, []);

  // Аутентификация пользователя
  const authenticateUser = async (userData: TelegramUser) => {
    try {
      setLoading(true);
      setError(null);
      
      // Получаем initData из Telegram WebApp
      const initData = (window as any).Telegram?.WebApp?.initData || 'dev';
      
      // Отправляем данные на бэкенд для аутентификации
      const response = await apiService.authenticateUser(userData, initData);
      
      // Сохраняем данные пользователя в состоянии
      setUser(response.user || userData);
      
      return response;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Authentication failed';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // Получение данных пользователя
  const getUserData = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await apiService.getUserData();
      return data;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to get user data';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // Сохранение данных пользователя
  const saveUserData = async (data: any) => {
    try {
      setLoading(true);
      setError(null);
      await apiService.saveUserData(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to save user data';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return {
    user,
    loading,
    error,
    authenticate: () => user ? authenticateUser(user) : Promise.reject(new Error('No user data')),
    getUserData,
    saveUserData,
  };
};

export default useApi;
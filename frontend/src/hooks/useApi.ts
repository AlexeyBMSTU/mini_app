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

  useEffect(() => {
    const initUser = async () => {
      try {
        try {
          const currentUser = await apiService.getCurrentUser();
          if (currentUser && currentUser.id) {
            setUser(currentUser);
            setLoading(false);
            return;
          }
        } catch {
          console.log('User not authenticated via cookies, proceeding with Telegram auth');
        }

        const telegramUser = telegramService.getUser();
        if (telegramUser) {
          setUser(telegramUser);
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

  const authenticateUser = async (userData: TelegramUser) => {
    try {
      setLoading(true);
      setError(null);
      
      const initData = (window as any).Telegram?.WebApp?.initData || 'dev';
      
      const response = await apiService.authenticateUser(userData, initData);
      
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
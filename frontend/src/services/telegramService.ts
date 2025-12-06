import { TelegramUser, TelegramWebApp } from '../types/telegram';

class TelegramService {
  private webApp: TelegramWebApp | null = null;

  init(): boolean {
    try {
      this.webApp = window.Telegram?.WebApp || null;
      
      if (this.webApp) {
        this.webApp.ready();
        this.webApp.expand();
        this.webApp.setHeaderColor('#ffffff');
        this.webApp.setBackgroundColor('#ffffff');
        return true;
      }
      return false;
    } catch (error) {
      console.error('Ошибка инициализации Telegram:', error);
      return false;
    }
  }

  isTelegram(): boolean {
    return !!this.webApp;
  }

  getUser(): TelegramUser | null {
    if (!this.webApp) return null;
    return this.webApp.initDataUnsafe?.user || null;
  }

  showAlert(message: string): void {
    if (this.webApp) {
      this.webApp.showAlert(message);
    } else {
      alert(message);
    }
  }

  close(): void {
    if (this.webApp) {
      this.webApp.close();
    }
  }

  getMockUser(): TelegramUser {
    return {
      id: 123456789,
      first_name: 'Разработчик',
      last_name: 'Тестовый',
      username: 'developer',
      language_code: 'ru',
      is_premium: true
    };
  }
}

export default new TelegramService();
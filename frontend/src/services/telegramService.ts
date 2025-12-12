import { TelegramUser, TelegramWebApp } from '../types/telegram'

class TelegramService {
	private webApp: TelegramWebApp | null = null
	private user: TelegramUser | undefined = undefined
  public isAvailable: boolean = false;

	init(): boolean {
		try {
			this.webApp = window.Telegram?.WebApp || null
			this.user = this.webApp?.initDataUnsafe?.user
      this.isAvailable = !!(this.webApp && this.user);

			if (this.webApp && this.isAvailable) {
				this.webApp.ready()
				this.webApp.expand()
				this.webApp.setHeaderColor('#ffffff')
				this.webApp.setBackgroundColor('#ffffff')
				return true
			}
			return false
		} catch (error) {
			console.error('Ошибка инициализации Telegram:', error)
			return false
		}
	}

	isTelegram(): boolean {
		return this.isAvailable;
	}

	getUser(): TelegramUser | null {
		if (!this.webApp || !this.isAvailable) return null
		return this.webApp.initDataUnsafe?.user || null
	}

	showAlert(message: string): void {
		if (this.webApp && this.isAvailable) {
			this.webApp.showAlert(message)
		}
	}

	close(): void {
		if (this.webApp && this.isAvailable) {
			this.webApp.close()
		}
	}

	getMockUser(): TelegramUser {
		return {
			id: 123456789,
			first_name: 'Разработчик',
			last_name: 'Тестовый',
			username: 'developer',
			language_code: 'ru',
			is_premium: true,
		}
	}
}

export default new TelegramService()

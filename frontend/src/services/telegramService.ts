import { TelegramUser, TelegramWebApp } from '../types/telegram'

class TelegramService {
  private webApp: TelegramWebApp | null = null
  private user: TelegramUser | undefined = undefined
  public isAvailable: boolean = false
  private backButtonCallback: (() => void) | null = null
  private mainButtonCallback: (() => void) | null = null

  init(): boolean {
    try {
      this.webApp = window.Telegram?.WebApp || null
      this.user = this.webApp?.initDataUnsafe?.user
      this.isAvailable = !!(this.webApp && this.user)

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
    return this.isAvailable
  }

  getUser(): TelegramUser | null {
    if (!this.webApp || !this.isAvailable)
      return {
        id: 123456789,
        first_name: 'Разработчик',
        last_name: 'Тестовый',
        username: 'developer',
        language_code: 'ru',
        is_premium: true,
      }
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

  showBackButton(callback?: () => void): void {
    if (this.webApp && this.isAvailable) {
      if (callback) {
        this.backButtonCallback = callback
        this.webApp.BackButton.onClick(callback)
      }
      this.webApp.BackButton.show()
    }
  }

  hideBackButton(): void {
    if (this.webApp && this.isAvailable) {
      if (this.backButtonCallback) {
        this.webApp.BackButton.offClick(this.backButtonCallback)
        this.backButtonCallback = null
      }
      this.webApp.BackButton.hide()
    }
  }

  disableMainButton(): void {
    if (this.webApp && this.isAvailable) {
      this.webApp.MainButton.disable()
    }
  }

  triggerHapticFeedback(type: 'light' | 'medium' | 'heavy' | 'rigid' | 'soft'): void {
    if (this.webApp && this.isAvailable) {
      this.webApp.HapticFeedback.impactOccurred(type)
    }
  }

  triggerNotificationFeedback(type: 'error' | 'success' | 'warning'): void {
    if (this.webApp && this.isAvailable) {
      this.webApp.HapticFeedback.notificationOccurred(type)
    }
  }

  triggerSelectionFeedback(): void {
    if (this.webApp && this.isAvailable) {
      this.webApp.HapticFeedback.selectionChanged()
    }
  }

  openLink(url: string, tryInstantView?: boolean): void {
    if (this.webApp && this.isAvailable) {
      this.webApp.openLink(url, { try_instant_view: tryInstantView || false })
    }
  }

  openTelegramLink(url: string): void {
    if (this.webApp && this.isAvailable) {
      this.webApp.openTelegramLink(url)
    }
  }
}

export default new TelegramService()

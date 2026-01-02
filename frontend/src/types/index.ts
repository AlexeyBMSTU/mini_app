export * from './telegram'
import { TelegramUser } from './telegram'

export interface AppState {
  user: TelegramUser | null
  loading: boolean
  isTelegram: boolean
}

export interface FeatureItem {
  title: string
  items: string[]
  available?: boolean
}

export interface ButtonConfig {
  text: string
  onClick: () => void
  style?: 'primary' | 'secondary' | 'telegram'
  disabled?: boolean
  visible?: boolean
}

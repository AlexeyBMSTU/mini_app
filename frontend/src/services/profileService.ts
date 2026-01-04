import apiService from './apiService'
import { TelegramUser } from '../types/telegram'

export interface ProfileData {
  email?: string
  phone?: string
  location?: string
  bio?: string
}

class ProfileService {
  async getUserProfile() {
    try {
      const response = await apiService.getCurrentUser()
      return response
    } catch (error) {
      console.error('Ошибка при получении данных профиля:', error)
      throw error
    }
  }

  async updateProfile(userData: Partial<TelegramUser>) {
    try {
      const response = await apiService.updateUser(userData)
      return response
    } catch (error) {
      console.error('Ошибка при обновлении профиля:', error)
      throw error
    }
  }

  async uploadAvatar(file: File) {
    try {
      const formData = new FormData()
      formData.append('avatar', file)

      const response = await fetch(`${apiService['baseUrl']}/api/user/avatar`, {
        method: 'POST',
        body: formData,
        credentials: 'include',
      })

      if (!response.ok) {
        throw new Error('Ошибка при загрузке аватара')
      }

      return await response.json()
    } catch (error) {
      console.error('Ошибка при загрузке аватара:', error)
      throw error
    }
  }

  async deleteAvatar() {
    try {
      const response = await fetch(`${apiService['baseUrl']}/api/user/avatar`, {
        method: 'DELETE',
        credentials: 'include',
      })

      if (!response.ok) {
        throw new Error('Ошибка при удалении аватара')
      }

      return await response.json()
    } catch (error) {
      console.error('Ошибка при удалении аватара:', error)
      throw error
    }
  }

  async getProfileData(): Promise<ProfileData> {
    try {
      const response = await apiService.getProfileData()
      return response
    } catch (error) {
      console.error('Ошибка при получении данных профиля:', error)
      throw error
    }
  }

  async updateProfileData(data: Partial<ProfileData>): Promise<ProfileData> {
    try {
      const response = await apiService.updateProfileData(data)
      return response
    } catch (error) {
      console.error('Ошибка при обновлении данных профиля:', error)
      throw error
    }
  }
}

export default new ProfileService()

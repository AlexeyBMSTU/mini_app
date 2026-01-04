import { DealType, PropertyType } from '../types/property'
import { TelegramUser } from '../types/telegram'

class ApiService {
  private baseUrl: string

  constructor() {
    this.baseUrl = 'http://localhost:8080'
  }

  async request(endpoint: string, options: RequestInit = {}) {
    const url = this.baseUrl + endpoint

    const defaultOptions: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    }

    const config = { ...defaultOptions, ...options }

    try {
      const response = await fetch(url, config)
      if (!response.ok) {
        console.error(response)
        throw new Error('HTTP error! status: ' + response.status)
      }

      return await response.json()
    } catch (error) {
      console.error('API request failed:', error)
      throw error
    }
  }

  async authenticateUser(userData: TelegramUser, initData: string) {
    return this.request('/api/auth/telegram/', {
      method: 'POST',
      body: JSON.stringify({
        user: userData,
        initData,
      }),
    })
  }

  async getCurrentUser() {
    return this.request('/api/user/me/')
  }

  async updateUser(userData: Partial<TelegramUser>) {
    return this.request('/api/user/me/', {
      method: 'PUT',
      body: JSON.stringify(userData),
    })
  }

  async getUserData() {
    return this.request('/api/user/data')
  }

  async saveUserData(data: any) {
    return this.request('/api/user/data', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async getToken() {
    return this.request('/api/auth/token', {
      method: 'POST',
    })
  }

  async getChats() {
    return this.request('/api/avito/messenger/chats/')
  }

  async getUserInfo() {
    return this.request('/api/avito/user/info/')
  }

  async getProperties(filters?: {
    propertyType?: PropertyType
    dealType?: DealType
    minPrice?: number
    maxPrice?: number
    minArea?: number
    maxArea?: number
    rooms?: number
  }) {
    const params = new window.URLSearchParams()

    if (filters) {
      if (filters.propertyType) params.append('type', filters.propertyType)
      if (filters.dealType) params.append('dealType', filters.dealType)
      if (filters.minPrice) params.append('minPrice', filters.minPrice.toString())
      if (filters.maxPrice) params.append('maxPrice', filters.maxPrice.toString())
      if (filters.minArea) params.append('minArea', filters.minArea.toString())
      if (filters.maxArea) params.append('maxArea', filters.maxArea.toString())
      if (filters.rooms) params.append('rooms', filters.rooms.toString())
    }

    const queryString = params.toString()
    return this.request('/api/properties' + (queryString ? '?' + queryString : ''))
  }

  async getPropertyById(id: string) {
    return this.request('/api/properties/' + id)
  }

  async createProperty(propertyData: any) {
    return this.request('/api/properties', {
      method: 'POST',
      body: JSON.stringify(propertyData),
    })
  }

  async getProfileData() {
    return this.request('/api/user/profile-data')
  }

  async updateProfileData(data: any) {
    return this.request('/api/user/profile-data', {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }
}

export default new ApiService()

import { TelegramUser } from '../types/telegram'

class ApiService {
	private baseUrl: string

	constructor() {
		this.baseUrl = `http://localhost:8080`
	}

	private async request(endpoint: string, options: RequestInit = {}) {
		const url = `${this.baseUrl}${endpoint}`

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
				throw new Error(`HTTP error! status: ${response.status}`)
			}

			return await response.json()
		} catch (error) {
			console.error('API request failed:', error)
			throw error
		}
	}

	async authenticateUser(userData: TelegramUser, initData: string) {
		return this.request('/api/auth/telegram', {
			method: 'POST',
			body: JSON.stringify({
				user: userData,
				initData,
			}),
		})
	}

	async getCurrentUser() {
		return this.request('/api/user/me')
	}

	async updateUser(userData: Partial<TelegramUser>) {
		return this.request('/api/user/me', {
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
			method: 'POST'
		})
	}
}

export default new ApiService()

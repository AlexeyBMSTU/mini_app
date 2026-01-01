class AvitoService {
	private baseUrl: string

	constructor() {
		this.baseUrl = `http://localhost:8080`
	}
	map = {
		cars: 9,
		motorcycles: 14,
		Trucks: 81,
		waterTranspor: 11,
		sparePartsAndAccessories: 10,
		apartments: 24,
		rooms: 23,
		houses: 25,
		landPlots: 26,
		garages: 85,
		commercialRealEstate: 42,
		realEstateAbroad: 86,
		residentialRentals: 338,
		jobOpportunities: 111,
		cvs: 112,
		serviceOfferings: 114,
	};

	private async request(endpoint: string, options: RequestInit = {}) {
		const url = `${this.baseUrl}${endpoint}`

		const defaultOptions: RequestInit = {
			headers: {
				'Content-Type': 'application/json',
			},
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

	async getItems() {
		try {
			const result = await this.request(`/api/avito/items/?category=${this.map.serviceOfferings}`, {
				method: 'GET',
			});
			return result;
		} catch (error) {
			console.error('Error in getItems:', error);
			throw error;
		}
	}

	async getChats() {
		try {
			const result = await this.request(`/api/avito/messenger/chats/`, {
				method: 'GET',
			});
			return result;
		} catch (error) {
			console.error('Error in getChats:', error);
			throw error;
		}
	}
}

export default new AvitoService();
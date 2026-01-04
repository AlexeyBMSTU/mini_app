import apiService from '@/services/apiService'
import { makeAutoObservable, runInAction } from 'mobx'
import { PropertyState, PropertyType, DealType } from '@/types/property'
import { mockProperties } from '@/mock/propertyMockData'

class BrowseStore {
  propertyState: PropertyState = {
    properties: [],
    loading: false,
    error: null,
    filters: {
      propertyType: undefined,
      dealType: undefined,
      minPrice: undefined,
      maxPrice: undefined,
      minArea: undefined,
      maxArea: undefined,
      rooms: undefined,
    },
  }

  useMockData: boolean = true

  constructor() {
    makeAutoObservable(this)
  }

  async fetchProperties() {
    try {
      runInAction(() => {
        this.propertyState.loading = true
        this.propertyState.error = null
      })

      if (this.useMockData) {
        await new Promise(resolve => window.setTimeout(resolve, 1000))

        runInAction(() => {
          this.propertyState.properties = mockProperties
          this.propertyState.loading = false
        })
      } else {
        const response = await apiService.getProperties(this.propertyState.filters)

        runInAction(() => {
          this.propertyState.properties = response.properties || []
          this.propertyState.loading = false
        })
      }
    } catch (error) {
      runInAction(() => {
        this.propertyState.error =
          error instanceof Error ? error.message : 'Failed to fetch properties'
        this.propertyState.loading = false
      })
    }
  }

  async fetchPropertyById(id: string) {
    try {
      runInAction(() => {
        this.propertyState.loading = true
        this.propertyState.error = null
      })

      if (this.useMockData) {
        await new Promise(resolve => window.setTimeout(resolve, 500))

        const property = mockProperties.find(p => p.id === id)
        if (!property) {
          throw new Error('Property not found')
        }

        runInAction(() => {
          const existingIndex = this.propertyState.properties.findIndex(p => p.id === id)
          if (existingIndex !== -1) {
            this.propertyState.properties[existingIndex] = property
          } else {
            this.propertyState.properties.push(property)
          }
          this.propertyState.loading = false
        })

        return property
      } else {
        const response = await apiService.getPropertyById(id)

        runInAction(() => {
          const existingIndex = this.propertyState.properties.findIndex(p => p.id === id)
          if (existingIndex !== -1) {
            this.propertyState.properties[existingIndex] = response
          } else {
            this.propertyState.properties.push(response)
          }
          this.propertyState.loading = false
        })

        return response
      }
    } catch (error) {
      runInAction(() => {
        this.propertyState.error =
          error instanceof Error ? error.message : 'Failed to fetch property'
        this.propertyState.loading = false
      })
      throw error
    }
  }

  setFilter<K extends keyof typeof this.propertyState.filters>(
    key: K,
    value: (typeof this.propertyState.filters)[K]
  ) {
    runInAction(() => {
      this.propertyState.filters[key] = value
    })
  }

  clearFilters() {
    runInAction(() => {
      this.propertyState.filters = {
        propertyType: undefined,
        dealType: undefined,
        minPrice: undefined,
        maxPrice: undefined,
        minArea: undefined,
        maxArea: undefined,
        rooms: undefined,
      }
    })
  }

  setPropertyType(propertyType: PropertyType) {
    this.setFilter('propertyType', propertyType)
  }

  setDealType(dealType: DealType) {
    this.setFilter('dealType', dealType)
  }

  setPriceRange(min: number | undefined, max: number | undefined) {
    this.setFilter('minPrice', min)
    this.setFilter('maxPrice', max)
  }

  setAreaRange(min: number | undefined, max: number | undefined) {
    this.setFilter('minArea', min)
    this.setFilter('maxArea', max)
  }

  setRoomsCount(rooms: number | undefined) {
    this.setFilter('rooms', rooms)
  }

  get properties() {
    return this.propertyState.properties
  }

  get loading() {
    return this.propertyState.loading
  }

  get error() {
    return this.propertyState.error
  }

  get filters() {
    return this.propertyState.filters
  }

  get filteredProperties() {
    const { properties, filters } = this.propertyState

    return properties.filter(property => {
      if (filters.propertyType && property.type !== filters.propertyType) return false
      if (filters.dealType && property.dealType !== filters.dealType) return false
      if (filters.minPrice && property.price < filters.minPrice) return false
      if (filters.maxPrice && property.price > filters.maxPrice) return false
      if (filters.minArea && property.area < filters.minArea) return false
      if (filters.maxArea && property.area > filters.maxArea) return false
      if (filters.rooms && property.rooms !== filters.rooms) return false

      return true
    })
  }

  get apartments() {
    return this.properties.filter(p => p.type === PropertyType.APARTMENT)
  }

  get houses() {
    return this.properties.filter(p => p.type === PropertyType.HOUSE)
  }

  get saleProperties() {
    return this.properties.filter(p => p.dealType === DealType.SALE)
  }

  get rentProperties() {
    return this.properties.filter(p => p.dealType === DealType.RENT)
  }
}

export const browseStore = new BrowseStore()

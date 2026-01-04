export interface PropertyImage {
  id: string
  url: string
  sizes: {
    '140x105': string
    '640x480': string
    '1280x960': string
  }
}

export interface PropertyOwner {
  id: number
  name: string
  phone?: string
  avatar?: {
    default: string
    images: {
      '128x128': string
      '192x192': string
      '24x24': string
      '256x256': string
      '36x36': string
      '48x48': string
      '64x64': string
      '72x72': string
      '96x96': string
    }
  }
}

export interface PropertyLocation {
  address: string
  city: string
  district?: string
  latitude: number
  longitude: number
}

export enum PropertyType {
  APARTMENT = 'apartment',
  HOUSE = 'house',
}

export enum DealType {
  SALE = 'sale',
  RENT = 'rent',
}

export interface Property {
  id: string
  title: string
  description: string
  type: PropertyType
  dealType: DealType
  price: number
  pricePerMeter?: number
  currency: string
  area: number
  rooms: number
  floor?: number
  totalFloors?: number
  yearBuilt?: number
  owner: PropertyOwner
  location: PropertyLocation
  images: PropertyImage[]
  features?: {
    hasBalcony?: boolean
    hasParking?: boolean
    hasElevator?: boolean
    hasFurniture?: boolean
    hasKitchen?: boolean
    hasInternet?: boolean
    hasAirConditioning?: boolean
    hasWashingMachine?: boolean
    hasDishwasher?: boolean
    hasTV?: boolean
    hasRefrigerator?: boolean
  }
  createdAt: string
  updatedAt: string
  isActive: boolean
}

export interface PropertyState {
  properties: Property[]
  loading: boolean
  error: string | null
  filters: {
    propertyType?: PropertyType
    dealType?: DealType
    minPrice?: number
    maxPrice?: number
    minArea?: number
    maxArea?: number
    rooms?: number
  }
}

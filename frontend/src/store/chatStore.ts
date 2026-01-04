import apiService from '@/services/apiService'
import { makeAutoObservable, runInAction } from 'mobx'

export interface Chat {
  id: string
  users: User[]
  created: number
  updated: number
  last_message: Message
  context: {
    type: string
    value: Item
  }
}

export interface User {
  id: number
  name: string
  public_user_profile: {
    avatar: {
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
    item_id: number
    url: string
    user_id: number
  }
}

export interface Message {
  id: string
  type: string
  created: number
  author_id: number
  direction: 'in' | 'out'
  read: boolean | number
  content: {
    text?: string
    image?: {
      sizes: ImageSizes
    }
    item?: {
      image_url: string
      item_url: string
      price_string: string
      title: string
    }
    link?: {
      text: string
      url: string
      preview: {
        title: string
        description: string
        domain: string
        url: string
        images: ImageSizes
      }
    }
    location?: {
      kind: string
      lat: number
      lon: number
      text: string
      title: string
    }
    call?: {
      status: string
      target_user_id: number
    }
    voice?: {
      voice_id: string
    }
    flow_id?: string
  }
}

export interface Item {
  id: number
  title: string
  price_string: string
  status_id: number
  url: string
  user_id: number
  images: {
    count: number
    main: {
      '140x105': string
    }
  }
}

export interface ImageSizes {
  '32x32': string
  '140x105': string
  '640x480': string
  '1280x960': string
}

export interface ChatState {
  chats: Chat[]
  loading: boolean
  error: string | null
}

class ChatStore {
  chatState: ChatState = {
    chats: [],
    loading: false,
    error: null,
  }

  constructor() {
    makeAutoObservable(this)
  }

  async fetchChats() {
    try {
      runInAction(() => {
        this.chatState.loading = true
        this.chatState.error = null
      })

      const response = await apiService.getChats()
      const responseUserInfo = await apiService.getUserInfo()

      runInAction(() => {
        this.chatState.chats = response.chats || []
        this.chatState.chats = this.chatState.chats.map(chat => ({
          ...chat,
          users: chat.users.filter(user => user.id !== responseUserInfo.id),
        }))

        this.chatState.loading = false
      })
    } catch (error) {
      runInAction(() => {
        this.chatState.error = error instanceof Error ? error.message : 'Failed to fetch chats'
        this.chatState.loading = false
      })
    }
  }

  get chats() {
    return this.chatState.chats
  }

  get loading() {
    return this.chatState.loading
  }

  get error() {
    return this.chatState.error
  }
}

export const chatStore = new ChatStore()

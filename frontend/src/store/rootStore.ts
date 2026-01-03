import { routeStore } from './routeStore'
import { chatStore } from './chatStore'
import { browseStore } from './browseStore'

export class RootStore {
  route = routeStore
  chat = chatStore
  browse = browseStore

  constructor() {}
}

export const rootStore = new RootStore()

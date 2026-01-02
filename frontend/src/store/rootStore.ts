import { routeStore } from './routeStore'
import { chatStore } from './chatStore'

export class RootStore {
  route = routeStore
  chat = chatStore

  constructor() {}
}

export const rootStore = new RootStore()

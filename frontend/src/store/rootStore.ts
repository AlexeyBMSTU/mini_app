import { routeStore } from './routeStore'

export class RootStore {
  route = routeStore

  constructor() {}
}

export const rootStore = new RootStore()

import { makeAutoObservable } from 'mobx'
import { NavigateFunction } from 'react-router-dom'

export interface RouteState {
  currentPath: string
  previousPath: string | null
  navigationHistory: string[]
}

class RouteStore {
  routeState: RouteState = {
    currentPath: '/',
    previousPath: null,
    navigationHistory: ['/'],
  }

  private navigateFn: NavigateFunction | null = null

  constructor() {
    makeAutoObservable(this)
  }

  setNavigateFunction(navigateFn: NavigateFunction) {
    this.navigateFn = navigateFn
  }

  navigate(path: string) {
    if (!this.navigateFn) {
      console.error(
        'Navigate function is not set. Make sure Router component is properly initialized.'
      )
      return
    }

    this.routeState.previousPath = this.routeState.currentPath
    this.routeState.currentPath = path

    if (!this.routeState.navigationHistory.includes(path)) {
      this.routeState.navigationHistory.push(path)
    }

    this.navigateFn(path)
  }

  goBack() {
    if (this.routeState.navigationHistory.length > 1 && this.navigateFn) {
      this.routeState.navigationHistory.pop()
      const previousPath =
        this.routeState.navigationHistory[this.routeState.navigationHistory.length - 1]
      this.routeState.previousPath = this.routeState.currentPath
      this.routeState.currentPath = previousPath

      this.navigateFn(-1)
    }
  }

  resetNavigation() {
    this.routeState = {
      currentPath: '/',
      previousPath: null,
      navigationHistory: ['/'],
    }

    if (this.navigateFn) {
      this.navigateFn('/')
    }
  }

  get currentPath() {
    return this.routeState.currentPath
  }

  get previousPath() {
    return this.routeState.previousPath
  }

  get canGoBack() {
    return this.routeState.navigationHistory.length > 1
  }
}

export const routeStore = new RouteStore()

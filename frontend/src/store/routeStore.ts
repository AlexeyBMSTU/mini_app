import { makeAutoObservable } from 'mobx'
import { NavigateFunction } from 'react-router-dom'
import telegramService from '../services/telegramService'

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
  private backButtonCallback: (() => void) | null = null

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

    this.navigateFn(path)
  }

  goBack() {
    if (this.canGoBack && this.navigateFn) {
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
    
    this.updateTelegramBackButton()
  }

  updateTelegramBackButton() {
    if (telegramService.isTelegram()) {
      if (this.backButtonCallback) {
        telegramService.hideBackButton()
      }

      this.backButtonCallback = () => {
        if (this.canGoBack) {
          this.goBack()
          telegramService.triggerHapticFeedback('light')
        } else {
          telegramService.showAlert('Вы на главной странице')
        }
      }
      console.log('alo')
      console.log(this.routeState)
      if (this.routeState.previousPath !== null) {
        telegramService.showBackButton(this.backButtonCallback)
      } else {
        telegramService.hideBackButton()
      }
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

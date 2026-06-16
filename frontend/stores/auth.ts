import { defineStore } from 'pinia'
import type { LoginResponse, User } from '~/types'

const TOKEN_KEY = 'budgetbooking.token'
const USER_KEY = 'budgetbooking.user'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: '' as string,
    user: null as User | null,
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
  },

  actions: {
    /** Load persisted session from localStorage. Call once on app start (client only). */
    init() {
      if (typeof window === 'undefined') return
      const token = window.localStorage.getItem(TOKEN_KEY)
      const userRaw = window.localStorage.getItem(USER_KEY)
      if (token) this.token = token
      if (userRaw) {
        try {
          this.user = JSON.parse(userRaw)
        } catch {
          this.user = null
        }
      }
    },

    setSession(token: string, user: User) {
      this.token = token
      this.user = user
      if (typeof window !== 'undefined') {
        window.localStorage.setItem(TOKEN_KEY, token)
        window.localStorage.setItem(USER_KEY, JSON.stringify(user))
      }
    },

    setLanguage(language: string) {
      if (this.user) this.user.language = language
      if (typeof window !== 'undefined' && this.user) {
        window.localStorage.setItem(USER_KEY, JSON.stringify(this.user))
      }
    },

    /** Authenticate with Google ID token, POST /user/login, store JWT + user + budgets. */
    async login(idToken: string) {
      const api = useApi()
      const res = await api<LoginResponse>('/user/login', {
        method: 'POST',
        body: { token: idToken },
        auth: false,
      })
      this.setSession(res.token, res.user)
      useBudgetsStore().setBudgets(res.budgets || [])
      return res.user
    },

    logout() {
      this.token = ''
      this.user = null
      if (typeof window !== 'undefined') {
        window.localStorage.removeItem(TOKEN_KEY)
        window.localStorage.removeItem(USER_KEY)
      }
      useBudgetsStore().clear()
    },
  },
})

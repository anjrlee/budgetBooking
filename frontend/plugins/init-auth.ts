import type { User, UserBudget } from '~/types'

export default defineNuxtPlugin(async () => {
  const authStore = useAuthStore()
  authStore.init()

  const budgetsStore = useBudgetsStore()
  if (authStore.token && !budgetsStore.loaded) {
    const api = useApi()
    try {
      const res = await api<{ user: User; budgets: UserBudget[] }>('/user/me')
      authStore.setSession(authStore.token, res.user)
      budgetsStore.setBudgets(res.budgets || [])
    } catch {
      // useApi already handles 401 by clearing the session and redirecting
    }
  }
})

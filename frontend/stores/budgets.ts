import { defineStore } from 'pinia'
import type { UserBudget } from '~/types'
import { SAVING_BUDGET_NAME } from '~/constants/budgetOptions'

export const useBudgetsStore = defineStore('budgets', {
  state: () => ({
    budgets: [] as UserBudget[],
    loaded: false,
  }),

  getters: {
    saving: (state) => state.budgets.find((b) => b.name === SAVING_BUDGET_NAME) || null,
    nonSaving: (state) => state.budgets.filter((b) => b.name !== SAVING_BUDGET_NAME),
    totalPercentage: (state) =>
      state.budgets.reduce((sum, b) => sum + (b.percentage || 0), 0),
    byId: (state) => (id: number | null) => state.budgets.find((b) => b.id === id) || null,
    usedIcons: (state) => (excludeId?: number) =>
      state.budgets.filter((b) => b.id !== excludeId).map((b) => b.icon),
    usedColors: (state) => (excludeId?: number) =>
      state.budgets.filter((b) => b.id !== excludeId).map((b) => b.color),
  },

  actions: {
    setBudgets(budgets: UserBudget[]) {
      this.budgets = budgets
      this.loaded = true
    },

    clear() {
      this.budgets = []
      this.loaded = false
    },
  },
})

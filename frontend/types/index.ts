export interface User {
  id: string
  email: string
  name: string
  picture?: string
  language: 'zh-Hant' | 'en' | string
}

export interface LoginResponse {
  token: string
  user: User
  budgets: UserBudget[]
}

export interface UserBudget {
  id: number
  name: string
  note?: string | null
  icon: string
  color: string
  percentage: number
  amount: number
}

export interface BudgetBooking {
  id: number
  budget_id: number | null
  budget_name: string | null
  is_income: boolean
  amount: number
  note?: string | null
  created_at: string
}

export interface AvailablePeriods {
  months: string[]
  years: number[]
}

export interface SpendingSummaryRow {
  budget_id: number | null
  budget_name: string
  icon?: string
  color?: string
  total_spent: number
  current_balance: number
}

export interface MonthlyTrendPoint {
  month: string
  total: number
}

export type DistributionMethod = 'auto' | 'specific'
export type BookingType = 'spending' | 'income'
export type ReportTab = 'spending' | 'income'
export type ReportScope = 'all' | 'specific'
export type ReportPeriod = 'monthly' | 'annual'

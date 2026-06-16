<script setup lang="ts">
import type { DataTableColumn } from '~/components/DataTable.vue'
import type {
  AvailablePeriods,
  BudgetBooking,
  MonthlyTrendPoint,
  ReportPeriod,
  ReportScope,
  ReportTab,
  SpendingSummaryRow,
  UserBudget,
} from '~/types'
import { PASTEL_COLORS } from '~/constants/budgetOptions'

const { t } = useI18n()
const api = useApi()
const budgetsStore = useBudgetsStore()
const { formatCurrency } = useCurrency()

const reportTab = ref<ReportTab>('spending')
const reportTabOptions = computed(() => [
  { value: 'spending', label: t('report.tabs.spending') },
  { value: 'income', label: t('report.tabs.income') },
])

const scope = ref<ReportScope>('all')
const scopeOptions = computed(() => [
  { value: 'all', label: t('report.scope.all') },
  { value: 'specific', label: t('report.scope.specific') },
])

const period = ref<ReportPeriod>('monthly')
const periodOptions = computed(() => [
  { value: 'monthly', label: t('report.period.monthly') },
  { value: 'annual', label: t('report.period.annual') },
])

const specificBudgetId = ref<number | null>(null)
const availablePeriods = ref<AvailablePeriods>({ months: [], years: [] })
const selectedMonth = ref<string | null>(null)
const selectedYear = ref<number | null>(null)

const loading = ref(false)
const summaryRows = ref<SpendingSummaryRow[]>([])
const transactions = ref<BudgetBooking[]>([])
const trendPoints = ref<MonthlyTrendPoint[]>([])
const incomeRows = ref<BudgetBooking[]>([])

const deleteTarget = ref<BudgetBooking | null>(null)
const deleting = ref(false)

const noPeriods = computed(() =>
  period.value === 'monthly' ? availablePeriods.value.months.length === 0 : availablePeriods.value.years.length === 0,
)

const summaryColumns: DataTableColumn[] = [
  { key: 'icon', label: '' },
  { key: 'budget_name', label: t('report.summary.budgetGroup') },
  { key: 'total_spent', label: t('report.summary.totalSpent'), align: 'right' },
  { key: 'current_balance', label: t('report.summary.currentBalance'), align: 'right' },
]

const transactionColumns = computed<DataTableColumn[]>(() => [
  { key: 'created_at', label: t('report.table.date') },
  { key: 'budget_name', label: t('report.table.budgetGroup') },
  { key: 'is_income', label: t('report.table.type') },
  { key: 'amount', label: t('report.table.amount'), align: 'right' },
  { key: 'note', label: t('report.table.note') },
  { key: 'actions', label: '', align: 'center' },
])

const pieLabels = computed(() => summaryRows.value.map((r) => r.budget_name))
const pieData = computed(() => summaryRows.value.map((r) => r.total_spent))
const pieColors = computed(() =>
  summaryRows.value.map((r, i) => r.color || budgetsStore.byId(r.budget_id)?.color || PASTEL_COLORS[i % PASTEL_COLORS.length]),
)

const visibleTransactions = computed(() => (reportTab.value === 'income' ? incomeRows.value : transactions.value))

const trendLabels = computed(() => trendPoints.value.map((p) => p.month))
const trendData = computed(() => trendPoints.value.map((p) => p.total))
const trendColor = computed(() => (specificBudgetId.value ? budgetsStore.byId(specificBudgetId.value)?.color : undefined))

function formatDate(value: string) {
  return new Date(value).toLocaleDateString()
}

async function loadAvailablePeriods() {
  const query: Record<string, unknown> = { type: reportTab.value }
  if (reportTab.value === 'spending' && scope.value === 'specific' && specificBudgetId.value) {
    query.budget_id = specificBudgetId.value
  }
  try {
    const res = await api<AvailablePeriods>('/report/availablePeriods', { query })
    availablePeriods.value = { months: res?.months || [], years: res?.years || [] }
  } catch {
    availablePeriods.value = { months: [], years: [] }
  }

  if (period.value === 'monthly') {
    selectedMonth.value = availablePeriods.value.months[availablePeriods.value.months.length - 1] || null
  } else {
    selectedYear.value = availablePeriods.value.years[availablePeriods.value.years.length - 1] || null
  }
}

async function loadReportData() {
  if (noPeriods.value) {
    summaryRows.value = []
    transactions.value = []
    trendPoints.value = []
    incomeRows.value = []
    return
  }

  loading.value = true
  try {
    if (reportTab.value === 'spending' && scope.value === 'all') {
      const query = period.value === 'monthly' ? { month: selectedMonth.value } : { year: selectedYear.value }
      const path = period.value === 'monthly' ? '/report/showSpendingReport' : '/report/showAnnualSpendingReport'
      const res = await api<{ summary: SpendingSummaryRow[]; transactions: BudgetBooking[] }>(path, { query })
      summaryRows.value = res.summary || []
      transactions.value = res.transactions || []
    } else if (reportTab.value === 'spending' && scope.value === 'specific') {
      const query: Record<string, unknown> = { budget_id: specificBudgetId.value }
      if (period.value === 'monthly') query.month = selectedMonth.value
      else query.year = selectedYear.value
      const path = period.value === 'monthly' ? '/report/showBudgetReport' : '/report/showAnnualBudgetReport'
      const res = await api<{ trend: MonthlyTrendPoint[]; transactions: BudgetBooking[] }>(path, { query })
      trendPoints.value = res.trend || []
      transactions.value = res.transactions || []
    } else {
      const query = period.value === 'monthly' ? { month: selectedMonth.value } : { year: selectedYear.value }
      const path = period.value === 'monthly' ? '/report/showIncomeReport' : '/report/showAnnualIncomeReport'
      const res = await api<{ rows: BudgetBooking[] }>(path, { query })
      incomeRows.value = res.rows || []
    }
  } catch {
    summaryRows.value = []
    transactions.value = []
    trendPoints.value = []
    incomeRows.value = []
  } finally {
    loading.value = false
  }
}

async function refreshAll() {
  await loadAvailablePeriods()
  await loadReportData()
}

watch(reportTab, () => {
  scope.value = 'all'
  specificBudgetId.value = budgetsStore.nonSaving[0]?.id ?? null
  refreshAll()
})

watch(scope, () => {
  if (scope.value === 'specific' && !specificBudgetId.value) {
    specificBudgetId.value = budgetsStore.nonSaving[0]?.id ?? null
  }
  refreshAll()
})

watch(specificBudgetId, () => {
  if (scope.value === 'specific') refreshAll()
})

watch(period, () => {
  refreshAll()
})

watch([selectedMonth, selectedYear], () => {
  loadReportData()
})

onMounted(() => {
  specificBudgetId.value = budgetsStore.nonSaving[0]?.id ?? null
  refreshAll()
})

function confirmDelete(row: BudgetBooking) {
  deleteTarget.value = row
}

async function deleteEntry() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    const res = await api<{ budgets: UserBudget[] }>('/booking/delete', {
      method: 'POST',
      body: { id: deleteTarget.value.id },
    })
    budgetsStore.setBudgets(res.budgets)
    deleteTarget.value = null
    await loadReportData()
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="mb-6 font-display text-2xl font-bold text-ink sm:text-3xl">{{ $t('report.title') }}</h1>

    <FunctionCard>
      <div class="mb-6 flex justify-center">
        <ToggleTabs v-model="reportTab" :options="reportTabOptions" />
      </div>

      <div class="mb-6 flex flex-wrap items-end gap-4">
        <div v-if="reportTab === 'spending'">
          <label class="field-label">{{ $t('report.scope.label') }}</label>
          <select v-model="scope" class="select-field">
            <option v-for="opt in scopeOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
          </select>
        </div>

        <div v-if="reportTab === 'spending' && scope === 'specific'">
          <label class="field-label">{{ $t('budget.selectGroup') }}</label>
          <select v-model.number="specificBudgetId" class="select-field">
            <option v-for="b in budgetsStore.nonSaving" :key="b.id" :value="b.id">{{ b.name }}</option>
          </select>
        </div>

        <div>
          <label class="field-label">{{ $t('report.period.label') }}</label>
          <select v-model="period" class="select-field">
            <option v-for="opt in periodOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
          </select>
        </div>

        <div v-if="!noPeriods && period === 'monthly'">
          <select v-model="selectedMonth" class="select-field">
            <option v-for="m in availablePeriods.months" :key="m" :value="m">{{ m }}</option>
          </select>
        </div>
        <div v-else-if="!noPeriods">
          <select v-model.number="selectedYear" class="select-field">
            <option v-for="y in availablePeriods.years" :key="y" :value="y">{{ y }}</option>
          </select>
        </div>
      </div>

      <p v-if="noPeriods" class="py-8 text-center text-sm text-ink-soft">{{ $t('report.noPeriods') }}</p>
      <div v-else-if="loading" class="py-8 text-center text-sm text-ink-soft">…</div>

      <template v-else>
        <!-- Spending: all categories -->
        <template v-if="reportTab === 'spending' && scope === 'all'">
          <PieChartCard v-if="pieData.length" :labels="pieLabels" :data="pieData" :colors="pieColors" class="mb-6" />

          <h3 class="mb-2 font-display text-lg font-bold text-ink">{{ $t('report.summary.title') }}</h3>
          <DataTable :columns="summaryColumns" :rows="summaryRows" row-key="budget_id" class="mb-6">
            <template #cell-icon="{ row }">
              <span
                class="flex h-8 w-8 items-center justify-center rounded-full"
                :style="{ backgroundColor: row.color || budgetsStore.byId(row.budget_id)?.color }"
              >
                <span class="material-symbols-outlined !text-[16px] text-ink">{{ row.icon || budgetsStore.byId(row.budget_id)?.icon }}</span>
              </span>
            </template>
            <template #cell-total_spent="{ value }">
              <span class="font-bold text-error">{{ formatCurrency(value) }}</span>
            </template>
            <template #cell-current_balance="{ value }">
              <span class="font-bold" :class="value < 0 ? 'text-error' : 'text-success'">{{ formatCurrency(value) }}</span>
            </template>
          </DataTable>
        </template>

        <!-- Spending: specific budget trend -->
        <template v-else-if="reportTab === 'spending' && scope === 'specific'">
          <LineChartCard v-if="trendLabels.length" :labels="trendLabels" :data="trendData" :color="trendColor" class="mb-6" />
        </template>

        <!-- Transactions -->
        <h3 class="mb-2 font-display text-lg font-bold text-ink">{{ $t('report.transactions') }}</h3>
        <BirdEmptyState v-if="!visibleTransactions.length" :message="$t('report.empty')" />
        <DataTable v-else :columns="transactionColumns" :rows="visibleTransactions" row-key="id">
          <template #cell-created_at="{ value }">{{ formatDate(value) }}</template>
          <template #cell-budget_name="{ value }">{{ value || $t('booking.distributed') }}</template>
          <template #cell-is_income="{ value }">{{ value ? $t('booking.income') : $t('booking.spending') }}</template>
          <template #cell-amount="{ row, value }">
            <span class="font-bold" :class="row.is_income ? 'text-success' : 'text-error'">
              {{ row.is_income ? '+' : '-' }}{{ formatCurrency(value) }}
            </span>
          </template>
          <template #cell-note="{ value }">
            <span class="text-ink-soft">{{ value || '—' }}</span>
          </template>
          <template #cell-actions="{ row }">
            <button type="button" class="icon-btn text-error" @click="confirmDelete(row)">
              <span class="material-symbols-outlined !text-[20px]">delete</span>
            </button>
          </template>
        </DataTable>
      </template>
    </FunctionCard>

    <ConfirmModal
      :open="!!deleteTarget"
      :title="$t('report.deleteConfirmTitle')"
      :body="$t('report.deleteConfirmBody')"
      danger
      @close="deleteTarget = null"
      @confirm="deleteEntry"
    />
  </div>
</template>

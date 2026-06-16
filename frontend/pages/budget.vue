<script setup lang="ts">
import type { MonthlyTrendPoint } from '~/types'

const api = useApi()
const budgetsStore = useBudgetsStore()
const { formatCurrency } = useCurrency()

const selectedBudgetId = ref<number | null>(budgetsStore.budgets[0]?.id ?? null)
const selectedBudget = computed(() => budgetsStore.byId(selectedBudgetId.value))

const loading = ref(false)
const trendLabels = ref<string[]>([])
const trendData = ref<number[]>([])

async function loadTrend() {
  if (!selectedBudgetId.value) return
  loading.value = true
  try {
    const now = new Date()
    const currentYear = now.getFullYear()
    const [curRes, prevRes] = await Promise.all([
      api<{ trend: MonthlyTrendPoint[] }>('/report/showAnnualBudgetReport', {
        query: { year: currentYear, budget_id: selectedBudgetId.value },
      }),
      api<{ trend: MonthlyTrendPoint[] }>('/report/showAnnualBudgetReport', {
        query: { year: currentYear - 1, budget_id: selectedBudgetId.value },
      }),
    ])
    const merged = [...(prevRes.trend || []), ...(curRes.trend || [])]
    const last12 = merged.slice(-12)
    trendLabels.value = last12.map((p) => p.month)
    trendData.value = last12.map((p) => p.total)
  } catch {
    trendLabels.value = []
    trendData.value = []
  } finally {
    loading.value = false
  }
}

watch(selectedBudgetId, loadTrend)
onMounted(loadTrend)
</script>

<template>
  <div>
    <h1 class="mb-6 font-display text-2xl font-bold text-ink sm:text-3xl">{{ $t('budget.title') }}</h1>

    <FunctionCard>
      <div class="mb-6">
        <label class="field-label">{{ $t('budget.selectGroup') }}</label>
        <select v-model.number="selectedBudgetId" class="select-field max-w-xs">
          <option v-for="b in budgetsStore.budgets" :key="b.id" :value="b.id">{{ b.name }}</option>
        </select>
      </div>

      <div v-if="selectedBudget" class="mb-8 flex items-center gap-4">
        <span
          class="flex h-16 w-16 flex-shrink-0 items-center justify-center rounded-full"
          :style="{ backgroundColor: selectedBudget.color }"
        >
          <span class="material-symbols-outlined !text-[32px] text-ink">{{ selectedBudget.icon }}</span>
        </span>
        <div>
          <p class="font-display text-xl font-bold text-ink">{{ selectedBudget.name }}</p>
          <p class="text-sm text-ink-soft">{{ $t('budget.currentBalance') }}</p>
          <p class="font-display text-2xl font-bold" :class="selectedBudget.amount < 0 ? 'text-error' : 'text-success'">
            {{ formatCurrency(selectedBudget.amount) }}
            <span v-if="selectedBudget.amount < 0" class="ml-1 text-sm font-bold">({{ $t('budget.overdraft') }})</span>
          </p>
        </div>
      </div>

      <h3 class="mb-2 font-display text-lg font-bold text-ink">{{ $t('budget.trend') }}</h3>
      <div v-if="loading" class="py-8 text-center text-sm text-ink-soft">…</div>
      <LineChartCard v-else-if="trendLabels.length" :labels="trendLabels" :data="trendData" :color="selectedBudget?.color" />
    </FunctionCard>
  </div>
</template>

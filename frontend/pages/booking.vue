<script setup lang="ts">
import type { DataTableColumn } from '~/components/DataTable.vue'
import type { BudgetBooking, UserBudget, BookingType, DistributionMethod } from '~/types'

const { t } = useI18n()
const api = useApi()
const budgetsStore = useBudgetsStore()
const { formatCurrency } = useCurrency()

const bookingType = ref<BookingType>('spending')
const bookingTypeOptions = computed(() => [
  { value: 'spending', label: t('booking.spending') },
  { value: 'income', label: t('booking.income') },
])

const distributionMethod = ref<DistributionMethod>('auto')
const distributionOptions = computed(() => [
  { value: 'auto', label: t('booking.form.auto') },
  { value: 'specific', label: t('booking.form.specific') },
])

const selectedBudgetId = ref<number | null>(null)
const amount = ref<number | null>(null)
const note = ref('')
const submitting = ref(false)
const feedback = ref<{ type: 'success' | 'error'; message: string } | null>(null)

watch(bookingType, () => {
  selectedBudgetId.value = null
  distributionMethod.value = 'auto'
  feedback.value = null
})

const canSubmit = computed(() => {
  if (!amount.value || amount.value <= 0) return false
  if (bookingType.value === 'spending') return !!selectedBudgetId.value
  if (distributionMethod.value === 'specific') return !!selectedBudgetId.value
  return true
})

async function submit() {
  if (!canSubmit.value || submitting.value) return
  submitting.value = true
  feedback.value = null
  try {
    const body: Record<string, unknown> = {
      is_income: bookingType.value === 'income',
      amount: Number(amount.value),
      note: note.value.trim() || null,
    }
    if (bookingType.value === 'spending' || distributionMethod.value === 'specific') {
      body.budget_id = selectedBudgetId.value
    }

    const res = await api<{ budgets: UserBudget[] }>('/booking/add', { method: 'POST', body })
    budgetsStore.setBudgets(res.budgets)
    feedback.value = { type: 'success', message: t('booking.success') }
    amount.value = null
    note.value = ''
    selectedBudgetId.value = null
    await loadRecent()
  } catch (err: any) {
    feedback.value = { type: 'error', message: err?.data?.message || t('booking.error') }
  } finally {
    submitting.value = false
  }
}

// Recent records
const recentIncome = ref<BudgetBooking[]>([])
const recentSpending = ref<BudgetBooking[]>([])
const loadingRecent = ref(false)
const deleteTarget = ref<BudgetBooking | null>(null)
const deleting = ref(false)

const tableColumns: DataTableColumn[] = [
  { key: 'created_at', label: t('booking.table.date') },
  { key: 'budget_name', label: t('booking.table.type') },
  { key: 'amount', label: t('booking.table.amount'), align: 'right' },
  { key: 'note', label: t('booking.table.note') },
  { key: 'actions', label: '', align: 'center' },
]

function currentMonth() {
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString()
}

async function loadRecentIncome() {
  loadingRecent.value = true
  try {
    const res = await api<{ rows: BudgetBooking[] }>('/report/showIncomeReport', {
      query: { month: currentMonth(), limit: 5 },
    })
    recentIncome.value = res.rows || []
  } catch {
    recentIncome.value = []
  } finally {
    loadingRecent.value = false
  }
}

async function loadRecentSpending() {
  loadingRecent.value = true
  try {
    const res = await api<{ transactions: BudgetBooking[] }>('/report/showSpendingReport', {
      query: { month: currentMonth(), limit: 5 },
    })
    recentSpending.value = res.transactions || []
  } catch {
    recentSpending.value = []
  } finally {
    loadingRecent.value = false
  }
}

async function loadRecent() {
  if (bookingType.value === 'income') await loadRecentIncome()
  else await loadRecentSpending()
}

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
    await loadRecent()
  } finally {
    deleting.value = false
  }
}

watch(bookingType, loadRecent)
onMounted(loadRecent)
</script>

<template>
  <div>
    <h1 class="mb-6 font-display text-2xl font-bold text-ink sm:text-3xl">{{ $t('booking.title') }}</h1>

    <FunctionCard>
      <div class="mb-6 flex justify-center">
        <ToggleTabs v-model="bookingType" :options="bookingTypeOptions" />
      </div>

      <div class="mx-auto max-w-md space-y-4">
        <div v-if="bookingType === 'spending'">
          <label class="field-label">{{ $t('booking.form.budgetGroup') }}</label>
          <select v-model.number="selectedBudgetId" class="select-field">
            <option :value="null" disabled>{{ $t('booking.form.selectBudgetGroup') }}</option>
            <option v-for="b in budgetsStore.budgets" :key="b.id" :value="b.id">{{ b.name }}</option>
          </select>
        </div>

        <div v-else>
          <label class="field-label">{{ $t('booking.form.distributionMethod') }}</label>
          <ToggleTabs v-model="distributionMethod" :options="distributionOptions" />

          <div v-if="distributionMethod === 'specific'" class="mt-3">
            <select v-model.number="selectedBudgetId" class="select-field">
              <option :value="null" disabled>{{ $t('booking.form.selectBudgetGroup') }}</option>
              <option v-for="b in budgetsStore.budgets" :key="b.id" :value="b.id">{{ b.name }}</option>
            </select>
          </div>
        </div>

        <div>
          <label class="field-label">{{ $t('booking.form.amount') }}</label>
          <input
            v-model.number="amount"
            type="number"
            min="0"
            step="0.01"
            class="input-field"
            :placeholder="$t('booking.form.amountPlaceholder')"
          >
        </div>

        <div>
          <label class="field-label">{{ $t('booking.form.note') }}</label>
          <textarea v-model="note" rows="2" class="input-field" :placeholder="$t('booking.form.note')" />
        </div>

        <button type="button" class="btn-primary w-full" :disabled="!canSubmit || submitting" @click="submit">
          {{ $t('booking.form.submit') }}
        </button>

        <BirdFeedback v-if="feedback" :type="feedback.type" :message="feedback.message" />
      </div>
    </FunctionCard>

    <FunctionCard
      :title="bookingType === 'income' ? $t('booking.recentIncome') : $t('booking.recentSpending')"
      class="mt-6"
    >
      <template v-if="bookingType === 'income'">
        <BirdEmptyState
          v-if="!loadingRecent && !recentIncome.length"
          image="/images/birds/bird-desk.png"
          :message="$t('booking.emptyIncome')"
        />
        <DataTable v-else :columns="tableColumns" :rows="recentIncome" row-key="id">
          <template #cell-created_at="{ value }">{{ formatDate(value) }}</template>
          <template #cell-budget_name="{ value }">{{ value || $t('booking.distributed') }}</template>
          <template #cell-amount="{ value }">
            <span class="font-bold text-success">+{{ formatCurrency(value) }}</span>
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

      <template v-else>
        <BirdEmptyState
          v-if="!loadingRecent && !recentSpending.length"
          image="/images/birds/bird-desk.png"
          :message="$t('booking.emptySpending')"
        />
        <DataTable v-else :columns="tableColumns" :rows="recentSpending" row-key="id">
          <template #cell-created_at="{ value }">{{ formatDate(value) }}</template>
          <template #cell-budget_name="{ value }">{{ value || '—' }}</template>
          <template #cell-amount="{ value }">
            <span class="font-bold text-error">-{{ formatCurrency(value) }}</span>
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
      :title="$t('booking.deleteConfirmTitle')"
      :body="$t('booking.deleteConfirmBody')"
      danger
      @close="deleteTarget = null"
      @confirm="deleteEntry"
    />
  </div>
</template>

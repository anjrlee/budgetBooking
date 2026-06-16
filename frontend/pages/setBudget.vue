<script setup lang="ts">
import type { DataTableColumn } from '~/components/DataTable.vue'
import type { UserBudget } from '~/types'
import { SAVING_BUDGET_NAME } from '~/constants/budgetOptions'

const { t } = useI18n()
const api = useApi()
const budgetsStore = useBudgetsStore()

onMounted(async () => {
  try {
    const res = await api<{ budgets: UserBudget[] }>('/user/me')
    budgetsStore.setBudgets(res.budgets || [])
  } catch {
    // keep showing whatever budgets are already in the store
  }
})

const showFormModal = ref(false)
const formMode = ref<'add' | 'edit'>('add')
const editingBudget = ref<UserBudget | null>(null)

const deleteTarget = ref<UserBudget | null>(null)
const showDeleteConfirm = computed(() => !!deleteTarget.value)
const deleting = ref(false)
const deleteError = ref('')

const pieLabels = computed(() => budgetsStore.budgets.map((b) => b.name))
const pieData = computed(() => budgetsStore.budgets.map((b) => b.percentage))
const pieColors = computed(() => budgetsStore.budgets.map((b) => b.color))

const columns: DataTableColumn[] = [
  { key: 'icon', label: '' },
  { key: 'name', label: t('setBudget.table.name') },
  { key: 'percentage', label: t('setBudget.table.percentage'), align: 'right' },
  { key: 'amount', label: t('setBudget.table.balance'), align: 'right' },
  { key: 'actions', label: t('setBudget.table.actions'), align: 'center' },
]

const tableRows = computed(() => budgetsStore.budgets)

const { formatCurrency } = useCurrency()

function openAddModal() {
  formMode.value = 'add'
  editingBudget.value = null
  showFormModal.value = true
}

function openEditModal(budget: UserBudget) {
  if (budget.name === SAVING_BUDGET_NAME) return
  formMode.value = 'edit'
  editingBudget.value = budget
  showFormModal.value = true
}

function closeFormModal() {
  showFormModal.value = false
}

function confirmDelete(budget: UserBudget) {
  if (budget.name === SAVING_BUDGET_NAME) return
  deleteError.value = ''
  deleteTarget.value = budget
}

async function deleteBudget() {
  if (!deleteTarget.value) return
  deleting.value = true
  deleteError.value = ''
  try {
    const res = await api<{ budgets: UserBudget[] }>('/setBudget/delete', {
      method: 'POST',
      body: { id: deleteTarget.value.id },
    })
    budgetsStore.setBudgets(res.budgets)
    deleteTarget.value = null
  } catch (err: any) {
    deleteError.value = err?.data?.message || t('booking.error')
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="mb-6 font-display text-2xl font-bold text-ink sm:text-3xl">{{ $t('setBudget.title') }}</h1>

    <FunctionCard>
      <div v-if="budgetsStore.budgets.length" class="flex flex-col items-center gap-6 sm:flex-row sm:items-center">
        <PieChartCard :labels="pieLabels" :data="pieData" :colors="pieColors" />
      </div>
      <BirdEmptyState v-else :message="$t('setBudget.hint')" />

      <div class="mt-4 flex justify-end">
        <button type="button" class="btn-primary" @click="openAddModal">
          <span class="material-symbols-outlined mr-1 align-middle !text-[18px]">add</span>
          {{ $t('setBudget.addGroup') }}
        </button>
      </div>
    </FunctionCard>

    <div v-if="budgetsStore.budgets.length" class="mt-6">
      <FunctionCard :title="$t('setBudget.groups')">
        <DataTable :columns="columns" :rows="tableRows" row-key="id">
          <template #cell-icon="{ row }">
            <span
              class="flex h-9 w-9 items-center justify-center rounded-full"
              :style="{ backgroundColor: row.color }"
            >
              <span class="material-symbols-outlined !text-[18px] text-ink">{{ row.icon }}</span>
            </span>
          </template>
          <template #cell-name="{ row }">
            <div>
              <p class="font-bold text-ink">{{ row.name }}</p>
              <p v-if="row.note" class="text-xs text-ink-soft">{{ row.note }}</p>
              <p v-if="row.name === 'saving'" class="text-xs text-ink-soft">{{ $t('setBudget.savingProtected') }}</p>
            </div>
          </template>
          <template #cell-percentage="{ value }">
            <span class="font-bold text-ink">{{ value }}%</span>
          </template>
          <template #cell-amount="{ value }">
            <span class="font-bold" :class="value < 0 ? 'text-error' : 'text-success'">{{ formatCurrency(value) }}</span>
          </template>
          <template #cell-actions="{ row }">
            <div class="flex items-center justify-center gap-1">
              <button
                v-if="row.name !== 'saving'"
                type="button"
                class="icon-btn"
                :aria-label="$t('setBudget.editGroup')"
                @click="openEditModal(row)"
              >
                <span class="material-symbols-outlined !text-[20px]">edit</span>
              </button>
              <button
                v-if="row.name !== 'saving'"
                type="button"
                class="icon-btn text-error"
                :aria-label="$t('setBudget.deleteGroup')"
                @click="confirmDelete(row)"
              >
                <span class="material-symbols-outlined !text-[20px]">delete</span>
              </button>
            </div>
          </template>
        </DataTable>
      </FunctionCard>
    </div>

    <BudgetFormModal
      :open="showFormModal"
      :mode="formMode"
      :budget="editingBudget"
      @close="closeFormModal"
    />

    <ConfirmModal
      :open="showDeleteConfirm"
      :title="$t('setBudget.deleteConfirmTitle')"
      :body="$t('setBudget.deleteConfirmBody', { name: deleteTarget?.name || '' })"
      :confirm-label="$t('setBudget.deleteGroup')"
      danger
      @close="deleteTarget = null"
      @confirm="deleteBudget"
    />
    <p v-if="deleteError" class="mt-2 text-sm font-bold text-error">{{ deleteError }}</p>
  </div>
</template>

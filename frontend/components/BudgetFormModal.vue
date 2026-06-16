<script setup lang="ts">
import { BUDGET_ICONS, PASTEL_COLORS } from '~/constants/budgetOptions'
import type { UserBudget } from '~/types'

const props = defineProps<{
  open: boolean
  mode: 'add' | 'edit'
  budget?: UserBudget | null
}>()

const emit = defineEmits<{ close: []; saved: [UserBudget[]] }>()

const { t } = useI18n()
const api = useApi()
const budgetsStore = useBudgetsStore()

const name = ref('')
const icon = ref('')
const color = ref('')
const note = ref('')
const newPercentage = ref<number | null>(null)
const percentageEdits = ref<{ id: number; name: string; icon: string; color: string; percentage: number }[]>([])
const submitting = ref(false)
const submitError = ref('')

function resetForm() {
  percentageEdits.value = budgetsStore.nonSaving.map((b) => ({
    id: b.id,
    name: b.name,
    icon: b.icon,
    color: b.color,
    percentage: b.percentage,
  }))

  if (props.mode === 'edit' && props.budget) {
    name.value = props.budget.name
    icon.value = props.budget.icon
    color.value = props.budget.color
    note.value = props.budget.note || ''
    newPercentage.value = null
  } else {
    name.value = ''
    note.value = ''
    newPercentage.value = null
    const usedI = budgetsStore.usedIcons()
    const usedC = budgetsStore.usedColors()
    icon.value = BUDGET_ICONS.find((i) => !usedI.includes(i)) || BUDGET_ICONS[0]
    color.value = PASTEL_COLORS.find((c) => !usedC.includes(c)) || PASTEL_COLORS[0]
  }
  submitError.value = ''
}

watch(
  () => props.open,
  (open) => {
    if (open) resetForm()
  },
)

const excludeId = computed(() => (props.mode === 'edit' ? props.budget?.id : undefined))
const usedIconsList = computed(() => budgetsStore.usedIcons(excludeId.value))
const usedColorsList = computed(() => budgetsStore.usedColors(excludeId.value))

const nonSavingTotal = computed(() => {
  const listTotal = percentageEdits.value.reduce((sum, b) => sum + (Number(b.percentage) || 0), 0)
  return listTotal + (props.mode === 'add' ? Number(newPercentage.value) || 0 : 0)
})
const savingPreview = computed(() => 100 - nonSavingTotal.value)
const isOverBudget = computed(() => savingPreview.value < 0)
const displayTotal = computed(() => (isOverBudget.value ? nonSavingTotal.value : 100))

const nameError = computed(() => {
  const trimmed = name.value.trim().toLowerCase()
  if (!trimmed) return ''
  const dup = budgetsStore.budgets.some((b) => b.name.toLowerCase() === trimmed && b.id !== excludeId.value)
  return dup ? t('setBudget.form.nameInUse') : ''
})

const canSubmit = computed(() => {
  if (props.mode === 'add' && (newPercentage.value === null || newPercentage.value <= 0)) return false
  return !!name.value.trim() && !nameError.value && !isOverBudget.value && !submitting.value
})

const title = computed(() => (props.mode === 'add' ? t('setBudget.addGroup') : t('setBudget.editGroup')))

async function submit() {
  if (!canSubmit.value) return
  submitting.value = true
  submitError.value = ''
  try {
    const updatedPercentages: Record<string, number> = {}
    percentageEdits.value.forEach((b) => {
      updatedPercentages[String(b.id)] = Number(b.percentage)
    })

    let res: { budgets: UserBudget[] }
    if (props.mode === 'add') {
      res = await api('/setBudget/add', {
        method: 'POST',
        body: {
          name: name.value.trim(),
          icon: icon.value,
          color: color.value,
          percentage: Number(newPercentage.value),
          note: note.value || null,
          updatedPercentages,
        },
      })
    } else {
      const original = props.budget!
      const percentageChanged = percentageEdits.value.some((b) => {
        const orig = budgetsStore.byId(b.id)
        return orig && Number(orig.percentage) !== Number(b.percentage)
      })
      res = await api('/setBudget/change', {
        method: 'POST',
        body: {
          id: original.id,
          name: name.value.trim(),
          icon: icon.value,
          color: color.value,
          note: note.value || null,
          ...(percentageChanged ? { updatedPercentages } : {}),
        },
      })
    }
    budgetsStore.setBudgets(res.budgets)
    emit('saved', res.budgets)
    emit('close')
  } catch (err: any) {
    submitError.value = err?.data?.message || t('booking.error')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <Modal :open="open" :title="title" @close="emit('close')">
    <div class="space-y-4">
      <div>
        <label class="field-label">{{ $t('setBudget.form.name') }}</label>
        <input v-model="name" type="text" class="input-field" :placeholder="$t('setBudget.form.namePlaceholder')">
        <p v-if="nameError" class="mt-1 text-xs font-bold text-error">{{ nameError }}</p>
      </div>

      <div>
        <label class="field-label">{{ $t('setBudget.form.icon') }}</label>
        <IconPicker v-model="icon" :disabled-icons="usedIconsList" />
      </div>

      <div>
        <label class="field-label">{{ $t('setBudget.form.color') }}</label>
        <ColorPicker v-model="color" :disabled-colors="usedColorsList" />
      </div>

      <div v-if="mode === 'add'">
        <label class="field-label">{{ $t('setBudget.form.percentage') }}</label>
        <input v-model.number="newPercentage" type="number" min="0" max="100" step="1" class="input-field">
      </div>

      <div>
        <label class="field-label">{{ $t('setBudget.form.note') }}</label>
        <textarea v-model="note" rows="2" class="input-field" :placeholder="$t('setBudget.form.notePlaceholder')" />
      </div>

      <div class="border-t border-outline/15 pt-4">
        <p class="field-label">{{ $t('setBudget.form.existingBudgets') }}</p>
        <ul class="space-y-2">
          <li v-for="b in percentageEdits" :key="b.id" class="flex items-center gap-3">
            <span
              class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full"
              :style="{ backgroundColor: b.color }"
            >
              <span class="material-symbols-outlined !text-[16px] text-ink">{{ b.icon }}</span>
            </span>
            <span class="flex-1 truncate text-sm font-bold text-ink">{{ b.name }}</span>
            <div class="flex items-center gap-1">
              <input
                v-model.number="b.percentage"
                type="number"
                min="0"
                max="100"
                step="1"
                class="input-field w-20 text-right"
              >
              <span class="text-sm text-ink-soft">%</span>
            </div>
          </li>
          <li v-if="budgetsStore.saving" class="flex items-center gap-3 opacity-70">
            <span
              class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full"
              :style="{ backgroundColor: budgetsStore.saving.color }"
            >
              <span class="material-symbols-outlined !text-[16px] text-ink">{{ budgetsStore.saving.icon }}</span>
            </span>
            <span class="flex-1 truncate text-sm font-bold text-ink">{{ budgetsStore.saving.name }}</span>
            <span class="text-sm text-ink-soft">
              {{ $t('setBudget.form.savingPreview', { from: budgetsStore.saving.percentage, to: Math.max(savingPreview, 0) }) }}
            </span>
          </li>
        </ul>
      </div>

      <div class="rounded-field bg-surface-container-low p-3 text-sm font-bold" :class="isOverBudget ? 'text-error' : 'text-success'">
        {{ isOverBudget ? $t('setBudget.form.total', { total: displayTotal }) : $t('setBudget.form.totalOk') }}
      </div>
      <p v-if="isOverBudget" class="text-xs font-bold text-error">{{ $t('setBudget.form.negativeSaving') }}</p>
      <p v-if="submitError" class="text-xs font-bold text-error">{{ submitError }}</p>

      <div class="flex justify-end gap-3 pt-2">
        <button type="button" class="btn-secondary" @click="emit('close')">{{ $t('common.cancel') }}</button>
        <button type="button" class="btn-primary" :disabled="!canSubmit" @click="submit">
          {{ $t('setBudget.form.submit') }}
        </button>
      </div>
    </div>
  </Modal>
</template>

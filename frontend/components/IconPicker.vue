<script setup lang="ts">
import { BUDGET_ICONS } from '~/constants/budgetOptions'

const props = defineProps<{
  modelValue: string
  disabledIcons?: string[]
}>()

const emit = defineEmits<{ 'update:modelValue': [string] }>()

function isDisabled(icon: string) {
  return icon !== props.modelValue && (props.disabledIcons || []).includes(icon)
}
</script>

<template>
  <div class="grid grid-cols-6 gap-2 sm:grid-cols-8">
    <button
      v-for="icon in BUDGET_ICONS"
      :key="icon"
      type="button"
      :disabled="isDisabled(icon)"
      class="flex h-10 w-10 items-center justify-center rounded-2xl border-2 transition"
      :class="[
        modelValue === icon ? 'border-primary bg-primary-light/40 text-primary-dark' : 'border-transparent bg-surface-container-low text-ink-soft',
        isDisabled(icon) ? 'cursor-not-allowed opacity-30' : 'hover:border-primary/50',
      ]"
      @click="emit('update:modelValue', icon)"
    >
      <span class="material-symbols-outlined !text-[20px]">{{ icon }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { PASTEL_COLORS } from '~/constants/budgetOptions'

const props = defineProps<{
  modelValue: string
  disabledColors?: string[]
}>()

const emit = defineEmits<{ 'update:modelValue': [string] }>()

function isDisabled(color: string) {
  return color !== props.modelValue && (props.disabledColors || []).includes(color)
}
</script>

<template>
  <div class="flex flex-wrap gap-2">
    <button
      v-for="color in PASTEL_COLORS"
      :key="color"
      type="button"
      :disabled="isDisabled(color)"
      class="h-9 w-9 rounded-full border-2 transition"
      :class="[
        modelValue === color ? 'border-primary scale-110' : 'border-transparent',
        isDisabled(color) ? 'cursor-not-allowed opacity-30' : 'hover:scale-110 saturate-150 brightness-95',
      ]"
      :style="{ backgroundColor: color }"
      @click="emit('update:modelValue', color)"
    />
  </div>
</template>
